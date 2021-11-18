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

type CoinGecko struct{}

type CoinGeckoMarketItem struct {
	ID           string    `json:"id"`
	Symbol       string    `json:"symbol"`
	Name         string    `json:"name"`
	CurrentPrice float32   `json:"current_price"`
	LastUpdated  time.Time `json:"last_updated"`
}

const coinGeckoGetMarketDataURL = "https://api.coingecko.com/api/v3/coins/markets"
const coinGeckoIDsQuery = "ids"
const coinGeckoVSCurrencyQuery = "vs_currency"

func (coinGecko CoinGecko) GetQuoteItems(ctx context.Context, targetCryptoIDs []string) ([]QuoteItem, error) {
	url, err := buildURLWithQueryParams(coinGeckoGetMarketDataURL, []query{
		{
			key:   coinGeckoIDsQuery,
			value: strings.Join(targetCryptoIDs, ","),
		},
		{
			key:   coinGeckoVSCurrencyQuery,
			value: defaultFiat,
		},
	})
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fail to request market data from CoinGecko")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jsonRes []CoinGeckoMarketItem
	if err := json.Unmarshal(body, &jsonRes); err != nil {
		return nil, err
	}

	var quoteItems []QuoteItem
	for _, v := range jsonRes {
		quoteItems = append(quoteItems, QuoteItem{
			Symbol:       strings.ToUpper(v.Symbol),
			Name:         v.Name,
			LastUpdated:  v.LastUpdated,
			BaseCurrency: defaultFiat,
			Price:        v.CurrentPrice,
		})
	}

	sortQuoteItemsAlphabeticallyASC(quoteItems)

	return quoteItems, nil
}
