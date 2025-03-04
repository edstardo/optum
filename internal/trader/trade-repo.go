package trader

import (
	"github.com/edstardo/mini-trader/pgk/postgres"
)

type TradeRepo interface {
	SaveTrade()
	GetUserTrades()
}

type tradesRepo struct {
	db *postgres.DB
}

func NewTradesRepo(db *postgres.DB) TradeRepo {
	return &tradesRepo{
		db: db,
	}
}

func (repo *tradesRepo) SaveTrade()     {}
func (repo *tradesRepo) GetUserTrades() {}
