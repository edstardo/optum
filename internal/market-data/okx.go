package marketdata

import (
	"context"
	"encoding/json"

	"github.com/edstardo/mini-trader/external/okx"
	"github.com/sirupsen/logrus"
)

type SourceOKX struct {
	okx *okx.OKX
}

func NewOKX() (PriceSource, error) {
	okx, err := okx.New()
	if err != nil {
		return nil, err
	}

	return &SourceOKX{
		okx: okx,
	}, nil
}

func (source *SourceOKX) GetPrices(ctx context.Context, tickers []string, msgChan chan Price) {
	// subscribe tickers
	source.okx.SubscribeTickers(tickers)

	// run okx msg reader goroutine
	go source.okx.ReadMessages(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case msgData := <-source.okx.MsgChan:
			logrus.Infof("raw okx data: %s", string(msgData))

			var msg okx.TickerDataMessage

			if err := json.Unmarshal(msgData, &msg); err != nil {
				logrus.Fatal(err)
			}

			if msg.Arg.Channel == okx.ChannelTickers && msg.Event == "" {
				msgChan <- &msg
			}
		}
	}
}
