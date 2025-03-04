package okx

import (
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type Arg struct {
	Channel string `json:"channel"`
	InstID  string `json:"instId"`
}

type SubscribeMessage struct {
	Op   string `json:"op"`
	Args []Arg  `json:"args"`
}

type TickerData struct {
	InstID string `json:"instId"`
	Last   string `json:"last"`
}

type Message struct {
	Event  string `json:"event"`
	Arg    Arg    `json:"arg"`
	Code   string `json:"code"`
	Msg    string `json:"msg"`
	ConnID string `json:"connId"`
}

type TickerDataMessage struct {
	Message
	Data []TickerData `json:"data"`
}

func (msg *TickerDataMessage) GetTicker() string {
	return msg.Arg.InstID
}

func (msg *TickerDataMessage) GetPrice() decimal.Decimal {
	price, err := decimal.NewFromString(msg.Data[0].Last)
	if err != nil {
		logrus.Error(err)
	}
	return price
}
