package oracle

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBuildURLWithQueryParams(t *testing.T) {
	queries := []query{
		{
			key:   "foo",
			value: "spam",
		},
		{
			key:   "bar",
			value: "ham",
		},
	}

	result, err := buildURLWithQueryParams("/hello", queries)
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, "/hello?bar=ham&foo=spam", result)
}

func TestSortQuoteItemsAlphabeticallyASC(t *testing.T) {
	itemA := QuoteItem{
		Symbol:       "A",
		Name:         "A Token",
		LastUpdated:  time.UnixMilli(1),
		BaseCurrency: "USD",
		Price:        1,
	}
	itemB := QuoteItem{
		Symbol:       "B",
		Name:         "B Token",
		LastUpdated:  time.UnixMilli(3),
		BaseCurrency: "USD",
		Price:        1,
	}
	itemC := QuoteItem{
		Symbol:       "C",
		Name:         "C Token",
		LastUpdated:  time.UnixMilli(2),
		BaseCurrency: "USD",
		Price:        1,
	}

	quoteItems := []QuoteItem{
		itemC,
		itemA,
		itemB,
	}

	sortQuoteItemsAlphabeticallyASC(quoteItems)

	assert.Equal(t, itemA, quoteItems[0])
	assert.Equal(t, itemB, quoteItems[1])
	assert.Equal(t, itemC, quoteItems[2])
}
