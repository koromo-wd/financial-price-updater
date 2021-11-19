package main

import (
	"testing"
	"time"

	"github.com/koromo-wd/priceupdater/oracle"
	"github.com/stretchr/testify/assert"
)

func TestCreateTradingPairs(t *testing.T) {
	quoteItems := []oracle.QuoteItem{
		{
			Symbol:       "A",
			Name:         "A Token",
			LastUpdated:  time.UnixMilli(1),
			BaseCurrency: "USD",
			Price:        1,
		},
		{
			Symbol:       "B",
			Name:         "B Token",
			LastUpdated:  time.UnixMilli(3),
			BaseCurrency: "USD",
			Price:        0.8,
		},
	}

	pairs := createTradingPairs(quoteItems)

	for i, pair := range pairs {
		quoteItem := quoteItems[i]
		assert.Equal(t, quoteItem.Symbol, pair.BaseSymbol)
		assert.Equal(t, "USD", pair.QuoteSymbol)
		assert.Equal(t, quoteItem.Price, pair.Price)
		assert.Equal(t, quoteItem.LastUpdated, pair.UpdatedTime)
	}
}
