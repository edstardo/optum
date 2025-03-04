package trader

import (
	"context"
	"encoding/json"

	natsio "github.com/nats-io/nats.go"
)

type Price interface {
	GetTicker() string
	GetPrice() string
}

type priceData struct {
	Ticker string `json:"ticker"`
	Price  string `json:"price"`
}

func (p *priceData) GetTicker() string {
	return p.Ticker
}
func (p *priceData) GetPrice() string {
	return p.Price
}

type MarketData interface {
	StreamPrices(ctx context.Context, priceChan chan Price)
}

type marketData struct {
	client *natsio.Conn
}

func NewMarketData(client *natsio.Conn) MarketData {
	return &marketData{
		client: client,
	}
}

func (m *marketData) StreamPrices(ctx context.Context, priceChan chan Price) {
	sub, err := m.client.Subscribe("prices.*", func(msg *natsio.Msg) {
		var priceData priceData

		if err := json.Unmarshal(msg.Data, &priceData); err != nil {
			return
		}
		priceChan <- &priceData
	})
	if err != nil {
		return
	}
	defer sub.Unsubscribe()

	<-ctx.Done()
}
