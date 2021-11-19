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

type fundPriceInfo struct {
	NavDate string  `json:"nav_date"`
	LastVal float32 `json:"last_val"`
}

const thb = "THB"
const bkkTz = "Asia/Bangkok"
const navDateOffSet = 10

const fundInfoURL = "https://api.sec.or.th/FundFactsheet/fund/class_fund"
const fundPriceURLTemplate = "https://api.sec.or.th/FundDailyInfo/%s/dailynav/%s"
const apiKeyHeader = "Ocp-Apim-Subscription-Key"
const navDateFormat = "2006-01-02"

var now = time.Now()

func (sec ThaiSec) GetQuoteItems(ctx context.Context, targetFundNames []string) ([]QuoteItem, error) {
	var quoteItems []QuoteItem

	for _, fundName := range targetFundNames {
		quoteItem, err := sec.getQuoteItem(ctx, fundName, getQueryNavDate(navDateOffSet))
		if err != nil {
			return nil, fmt.Errorf("fundName=%s %w", fundName, err)
		}

		quoteItems = append(quoteItems, *quoteItem)
	}

	sortQuoteItemsAlphabeticallyASC(quoteItems)

	return quoteItems, nil
}

func (sec ThaiSec) getQuoteItem(ctx context.Context, fundName, queryNavDate string) (*QuoteItem, error) {
	fundInfo, err := sec.getFundInfo(ctx, fundName)
	if err != nil {
		return nil, fmt.Errorf("fail to get fund info from Thai SEC API: %w", err)
	}

	fundPrice, err := sec.getFundPrice(ctx, fundInfo.ProjectID, queryNavDate)
	if err != nil {
		return nil, fmt.Errorf("fundID=%s queryNavDate=%s fail to request fund price from Thai SEC API: %w", fundInfo.ProjectID, queryNavDate, err)
	}

	timeLoc, err := getTimeLoc(bkkTz)
	if err != nil {
		return nil, err
	}

	parsedTime, err := time.ParseInLocation(navDateFormat, fundPrice.NavDate, timeLoc)
	if err != nil {
		return nil, err
	}

	return &QuoteItem{
		Symbol:       fundName,
		Name:         fundInfo.ProjectABBRName,
		LastUpdated:  parsedTime,
		BaseCurrency: thb,
		Price:        fundPrice.LastVal,
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
		return nil, fmt.Errorf("request returns statusCode=%d", resp.StatusCode)
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

func (sec ThaiSec) getFundPrice(ctx context.Context, fundID, queryNavDate string) (*fundPriceInfo, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(fundPriceURLTemplate, fundID, queryNavDate), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set(apiKeyHeader, sec.FundDailyInfoAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request returns statusCode=%d", resp.StatusCode)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jsonRes fundPriceInfo
	if err := json.Unmarshal(body, &jsonRes); err != nil {
		return nil, err
	}

	return &jsonRes, nil
}

func getQueryNavDate(pastDayOffset int) string {
	return now.AddDate(0, 0, -1*pastDayOffset).Format(navDateFormat)
}

func getTimeLoc(countryTz string) (*time.Location, error) {
	loc, err := time.LoadLocation(countryTz)
	if err != nil {
		return nil, err
	}
	return loc, nil
}
