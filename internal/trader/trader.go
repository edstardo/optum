package trader

import (
	"context"

	"github.com/sirupsen/logrus"
)

type service struct {
	marketData MarketData
	memory     Memory
	priceChan  chan Price
	tradesRepo TradeRepo
}

func New(marketData MarketData, memory Memory, tradesRepo TradeRepo) *service {
	return &service{
		marketData: marketData,
		memory:     memory,
		priceChan:  make(chan Price),
		tradesRepo: tradesRepo,
	}
}

func (s *service) GetAndSavePrices(ctx context.Context) {
	go s.marketData.StreamPrices(ctx, s.priceChan)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case price := <-s.priceChan:
				logrus.Infof("%v: %v", price.GetTicker(), price.GetPrice())

				if err := s.memory.SavePrice(ctx, price.GetTicker(), price.GetPrice()); err != nil {
					logrus.Error(err)
				}
			}
		}
	}()
}
