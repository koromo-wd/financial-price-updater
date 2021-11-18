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

const cmcQuoteURL string = "https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest"
const cmcAPIKeyQuery string = "CMC_PRO_API_KEY"
const cmcSymbolQuery string = "symbol"

type CMC struct {
	APIKey string
}

type CMCQuoteJSONResponse struct {
	Status map[string]interface{}  `json:"status"`
	Data   map[string]CMCQuoteItem `json:"data"`
}

type CMCQuoteItem struct {
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
	url, err := buildURLWithQueryParams(cmcQuoteURL, []query{
		{
			key:   cmcAPIKeyQuery,
			value: cmc.APIKey,
		},
		{
			key:   cmcSymbolQuery,
			value: strings.Join(targetCryptoSymbols, ","),
		},
	})
	if err != nil {
		return nil, err
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

	var jsonRes CMCQuoteJSONResponse
	if err := json.Unmarshal(body, &jsonRes); err != nil {
		return nil, err
	}

	var quoteItems []QuoteItem
	for _, v := range jsonRes.Data {
		quoteItems = append(quoteItems, QuoteItem{
			Symbol:       v.Symbol,
			Name:         v.Name,
			LastUpdated:  v.LastUpdated,
			BaseCurrency: defaultFiat,
			Price:        v.Quote.USD.Price,
		})
	}

	sortQuoteItemsAlphabeticallyASC(quoteItems)

	return quoteItems, nil
}
