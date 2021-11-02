package oracle

import (
	"context"
	"net/url"
	"sort"
	"time"
)

type CryptoOracle interface {
	GetQuoteItems(ctx context.Context, targetCrypto []string) ([]QuoteItem, error)
}

type QuoteItem struct {
	Symbol      string
	Name        string
	Slug        string
	LastUpdated time.Time
	USDPrice    float32
}

type query struct {
	key   string
	value string
}

func buildURLWithQueryParams(baseURL string, queries []query) (string, error) {
	url, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	params := url.Query()
	for _, query := range queries {
		params.Add(query.key, query.value)
	}
	url.RawQuery = params.Encode()

	return url.String(), nil
}

// sortQuoteItems (alphabetically asc)
func sortQuoteItems(quoteItems []QuoteItem) {
	sort.Slice(quoteItems, func(i, j int) bool {
		return quoteItems[i].Symbol < quoteItems[j].Symbol
	})
}
