package trader

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type Quote struct {
	QuoteID string `json:"quote_id"`
	UserID  string `json:"user_id"`

	Ticker string `json:"ticker"`

	Quantity decimal.Decimal `json:"quantity"`
	Price    decimal.Decimal `json:"price"`
	Amount   decimal.Decimal `json:"amount"`

	Side string `json:"side"`

	CreatedAt time.Time `json:"created_at"`
}

func (q *Quote) IsValid() (bool, error) {
	if time.Now().After(q.CreatedAt.Add(DefaultTTL * time.Second)) {
		return false, fmt.Errorf("quote has already expired")
	}
	return true, nil
}
