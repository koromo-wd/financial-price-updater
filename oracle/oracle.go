package oracle

import (
	"context"
	"net/url"
)

type CryptoOracle interface {
	GetQuoteItems(ctx context.Context, targetCryptoSymbols []string) []QuoteItem
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
