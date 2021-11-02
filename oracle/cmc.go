package oracle

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const apiKeyQuery string = "CMC_PRO_API_KEY"
const symbolQuery string = "symbol"
const quoteURL string = "https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest"

type CMC struct {
	APIKey string
}

type QuoteJSONResponse struct {
	Status map[string]interface{} `json:"status"`
	Data   map[string]QuoteItem   `json:"data"`
}

type QuoteItem struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Symbol      string    `json:"symbol"`
	Slug        string    `json:"slug"`
	LastUpdated time.Time `json:"last_updated"`
	Quote       struct {
		USD struct {
			Price float32 `json:"price"`
		} `json:"USD"`
	} `json:"quote"`
}

func (cmc CMC) GetQuoteItems(ctx context.Context, targetCryptoSymbols []string) ([]QuoteItem, error) {
	url, err := buildURLWithQueryParams(quoteURL, []query{
		{
			key:   apiKeyQuery,
			value: cmc.APIKey,
		},
		{
			key:   symbolQuery,
			value: strings.Join(targetCryptoSymbols, ","),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("fail to build quote URL")
	}
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fail to request quote data from CoinMarketCap")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jsonRes QuoteJSONResponse
	if err := json.Unmarshal(body, &jsonRes); err != nil {
		return nil, err
	}

	var targetItems []QuoteItem

	for _, v := range jsonRes.Data {
		targetItems = append(targetItems, v)
	}

	return targetItems, nil
}
