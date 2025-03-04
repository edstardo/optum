package trader

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type Memory interface {
	SavePrice(ctx context.Context, ticker string, price string) error
	GetPrice(ctx context.Context, ticker string) (string, error)

	SaveQuote(ctx context.Context, quote Quote) error
	GetQuote(ctx context.Context, userId, quoteID string) (*Quote, error)
}

func NewMemory(client *redis.Client) Memory {
	return &memory{
		client: client,
	}
}

type memory struct {
	client *redis.Client
}

func (m *memory) SavePrice(ctx context.Context, ticker string, price string) error {
	key := fmt.Sprintf("price-%s", ticker)
	ttl := 5 * time.Second

	cmd := m.client.Set(ctx, key, price, ttl)
	if err := cmd.Err(); err != nil {
		return err
	}

	return nil
}

func (m *memory) GetPrice(ctx context.Context, ticker string) (string, error) {
	key := fmt.Sprintf("price-%s", ticker)

	cmd := m.client.Get(ctx, key)
	if err := cmd.Err(); err != nil {
		return "", err
	}

	return cmd.Val(), nil
}

func (m *memory) SaveQuote(ctx context.Context, quote Quote) error {
	key := fmt.Sprintf("quote-%s-%s", quote.UserID, quote.QuoteID)

	b, err := json.Marshal(quote)
	if err != nil {
		logrus.Error(err)
		return err
	}

	cmd := m.client.Set(ctx, key, b, 1*time.Minute)
	if err := cmd.Err(); err != nil {
		return err
	}
	return nil
}

func (m *memory) GetQuote(ctx context.Context, userID, quoteID string) (*Quote, error) {
	key := fmt.Sprintf("quote-%s-%s", userID, quoteID)

	cmd := m.client.Get(ctx, key)
	if err := cmd.Err(); err != nil {
		logrus.Error(err)
		return nil, err
	}

	b, _ := cmd.Bytes()

	var quote Quote
	if err := json.Unmarshal(b, &quote); err != nil {
		logrus.Error(err)
		return nil, err
	}

	return &quote, nil
}
