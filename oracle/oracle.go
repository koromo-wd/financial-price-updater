package oracle

import (
	"context"
	"net/url"
	"sort"
	"time"
)

type Oracle interface {
	GetQuoteItems(ctx context.Context, queryTargets []string) ([]QuoteItem, error)
}

type QuoteItem struct {
	Symbol       string
	Name         string
	LastUpdated  time.Time
	BaseCurrency string
	Price        float32
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

func sortQuoteItemsAlphabeticallyASC(quoteItems []QuoteItem) {
	sort.Slice(quoteItems, func(i, j int) bool {
		return quoteItems[i].Symbol < quoteItems[j].Symbol
	})
}
