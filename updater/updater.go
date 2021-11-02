package updater

import (
	"context"
	"time"
)

type Updater interface {
	UpdatePrice(ctx context.Context, tradingPairs []TradingPair) error
}

type TradingPair struct {
	BaseSymbol  string
	QuoteSymbol string
	Price       float32
	UpdatedTime time.Time
}
