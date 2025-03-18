package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"

	marketdata "github.com/edstardo/optum/internal/market-data"
	"github.com/sirupsen/logrus"
)

const (
	NATS_SERVER = "nats://localhost:4222,nats://localhost:4223,nats://localhost:4224"
	OKX_TICKERS = "BTC-USDT,ETH-USDT"
)

func main() {
	natsServer := os.Getenv("NATS_SERVER")
	tickers := strings.Split(OKX_TICKERS, ",")

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-exit
		cancel()
	}()

	okx, err := marketdata.NewPriceSource(marketdata.SourceNameOKX)
	if err != nil {
		logrus.Fatal(err)
	}

	nats, err := marketdata.NewNats(natsServer)
	if err != nil {
		logrus.Fatal(err)
	}

	marketDataSvc := marketdata.New(okx, nats)

	marketDataSvc.GetPrices(ctx, tickers)

	logrus.Info("shutdown gracefully")
}
