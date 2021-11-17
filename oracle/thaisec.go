package oracle

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type ThaiSec struct {
	FundFactAPIKey      string
	FundDailyInfoAPIKey string
}

type fundInfo struct {
	ProjectID       string `json:"proj_id"`
	ProjectABBRName string `json:"proj_abbr_name"`
}

type fundPriceResponse struct {
	LastUpdatedDate string  `json:"last_upd_date"`
	NavDate         string  `json:"nav_date"`
	LastVal         float32 `json:"last_val"`
}

const thb = "THB"
const bkkTz = "Asia/Bangkok"
const navDateOffSet = 5

const fundInfoURL = "https://api.sec.or.th/FundFactsheet/fund/class_fund"
const fundPriceURLTemplate = "https://api.sec.or.th/FundDailyInfo/%s/dailynav/%s"
const apiKeyHeader = "Ocp-Apim-Subscription-Key"
const navDateFormat = "2006-01-02"

func (sec ThaiSec) GetQuoteItems(ctx context.Context, targetFundNames []string) ([]QuoteItem, error) {
	var quoteItems []QuoteItem

	for _, fundName := range targetFundNames {
		quoteItem, err := sec.getQuoteItem(ctx, fundName, getQueryNavDate(navDateOffSet))
		if err != nil {
			return nil, err
		}

		quoteItems = append(quoteItems, *quoteItem)
	}

	sortQuoteItems(quoteItems)

	return quoteItems, nil
}

func (sec ThaiSec) getQuoteItem(ctx context.Context, fundName, queryNavDate string) (*QuoteItem, error) {
	fundInfo, err := sec.getFundInfo(ctx, fundName)
	if err != nil {
		return nil, fmt.Errorf("fail to get TH Sec Fund Info: %w", err)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf(fundPriceURLTemplate, fundInfo.ProjectID, queryNavDate), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set(apiKeyHeader, sec.FundDailyInfoAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fail to request fund price of %s from Thai Sec API", fundName)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fail to request fund price of %s with status: %d", fundName, resp.StatusCode)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jsonRes fundPriceResponse
	if err := json.Unmarshal(body, &jsonRes); err != nil {
		return nil, err
	}

	timeLoc, err := getTimeLoc(bkkTz)
	if err != nil {
		return nil, err
	}

	parsedTime, err := time.ParseInLocation(navDateFormat, jsonRes.NavDate, timeLoc)
	if err != nil {
		return nil, err
	}

	return &QuoteItem{
		Symbol:      fundName,
		Name:        fundInfo.ProjectABBRName,
		Slug:        "",
		LastUpdated: parsedTime,
		Price: map[string]float32{
			thb: jsonRes.LastVal,
		},
	}, nil
}

func (sec ThaiSec) getFundInfo(ctx context.Context, fundName string) (*fundInfo, error) {
	reqBody, err := json.Marshal(map[string]string{
		"name": fundName,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fundInfoURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set(apiKeyHeader, sec.FundFactAPIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fail to request fund info with status: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jsonRes []fundInfo
	if err := json.Unmarshal(respBody, &jsonRes); err != nil {
		return nil, err
	}

	return &jsonRes[0], nil
}

func getQueryNavDate(pastDayOffset int) string {
	return time.Now().AddDate(0, 0, -1*pastDayOffset).Format(navDateFormat)
}

func getTimeLoc(countryTz string) (*time.Location, error) {
	loc, err := time.LoadLocation(countryTz)
	if err != nil {
		return nil, err
	}
	return loc, nil
}
