package marketdata

import (
	"context"
	"encoding/json"

	"github.com/shopspring/decimal"
)

type Service struct {
	publisher   PricePublisher
	source      PriceSource // external price source, in this case just OKX but can be anything
	priceInChan chan Price  // used for receiving prices that implements Price interface from source
}

type PriceData struct {
	Ticker string `json:"ticker"`
	Price  string `json:"price"`
}

func New(source PriceSource, publisher PricePublisher) *Service {
	return &Service{
		source:      source,
		publisher:   publisher,
		priceInChan: make(chan Price),
	}
}

func (s *Service) GetPrices(ctx context.Context, tickers []string) {
	go s.source.GetPrices(ctx, tickers, s.priceInChan)

	for {
		select {
		case <-ctx.Done():
			return
		case price := <-s.priceInChan:
			s.PublishPrices(ctx, price.GetTicker(), price.GetPrice())
		}
	}
}

func (s *Service) PublishPrices(ctx context.Context, ticker string, price decimal.Decimal) {
	priceData := PriceData{
		Ticker: ticker,
		Price:  price.StringFixed(8),
	}

	b, err := json.Marshal(&priceData)
	if err != nil {
		return
	}

	_ = s.publisher.Publish(ticker, b)
}
