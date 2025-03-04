package okx

import (
	"context"
	"log"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

const (
	OpSubscribe = "subscribe"

	EventSubscribe = "subscribe"
	EventError     = "error"

	ChannelTickers = "tickers"

	PublicBaseURL = "wss://ws.okx.com:8443/ws/v5/public"
)

type OKX struct {
	conn    *websocket.Conn
	MsgChan chan []byte
}

func New() (*OKX, error) {
	conn, _, err := websocket.DefaultDialer.Dial(PublicBaseURL, nil)
	if err != nil {
		return nil, err
	}

	return &OKX{
		conn:    conn,
		MsgChan: make(chan []byte),
	}, nil
}

func (o *OKX) SubscribeTickers(tickers []string) error {
	var args []Arg

	for _, ticker := range tickers {
		args = append(args, Arg{
			Channel: ChannelTickers,
			InstID:  ticker,
		})
	}

	subMsg := SubscribeMessage{
		Op:   OpSubscribe,
		Args: args,
	}

	err := o.conn.WriteJSON(&subMsg)
	if err != nil {
		return err
	}

	return nil
}

func (o *OKX) ReadMessages(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			logrus.Info("okx message reader exited")
			return
		default:
			_, message, err := o.conn.ReadMessage()
			if err != nil {
				log.Fatal(err)
				return
			}
			o.MsgChan <- message
		}
	}
}

func (o *OKX) Close() error {
	return o.conn.Close()
}
