package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/edstardo/optum/internal/trader"
	"github.com/edstardo/optum/pgk/postgres"
	natsio "github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func main() {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		<-exit
		cancel()
	}()

	// init nats client
	natsClient, err := natsio.Connect(os.Getenv("NATS_SERVER"))
	if err != nil {
		logrus.Fatal(err)
	}

	marketData := trader.NewMarketData(natsClient)

	// init redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_SERVER"),
	})

	// price repo
	memory := trader.NewMemory(redisClient)

	// init postgres
	pgdb, err := postgres.NewPostgresDB(os.Getenv("POSTGRES_URL"))
	if err != nil {
		logrus.Fatal(err)
	}

	// init trades repo
	tradesRepo := trader.NewTradesRepo(pgdb)

	// init trader services
	traderSvc := trader.New(marketData, memory, tradesRepo)

	// setup api router
	router := trader.NewRouter(traderSvc)

	// setup api server
	httpServer := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatal(err)
		}
	}()

	// get and save prices to repo
	traderSvc.GetAndSavePrices(ctx)

	<-ctx.Done()

	// gracefully shutdown http server
	sCtx, sCancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	defer sCancel()

	if err := httpServer.Shutdown(sCtx); err != nil {
		logrus.Error(err)
	}

	logrus.Info("gracefully shutdown")
}
