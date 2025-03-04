package marketdata

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"
)

const (
	SourceNameOKX = "okx"
)

// PriceSource is an iterface that should be implemented by all price sources
type PriceSource interface {
	GetPrices(ctx context.Context, tickers []string, priceChan chan Price)
}

// PriceData is an iterface that should be implemented by price source data
type Price interface {
	GetTicker() string
	GetPrice() decimal.Decimal
}

// NewPriceSource create a source instance.
// This implements a builder pattern that resturn a `Source` interface
// which can be added with more concrete implementations of this interface.
func NewPriceSource(sourceName string) (PriceSource, error) {
	// builder pattern

	if sourceName == SourceNameOKX {
		return NewOKX()
	}

	return nil, fmt.Errorf("unknown price source name")
}
