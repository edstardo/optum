package marketdata

import (
	"fmt"

	natsio "github.com/nats-io/nats.go"
)

type Nats struct {
	client *natsio.Conn
}

func NewNats(url string) (*Nats, error) {
	nc, err := natsio.Connect(url)
	if err != nil {
		return nil, err
	}

	return &Nats{
		client: nc,
	}, nil
}

func (n *Nats) Publish(topic string, data []byte) error {
	subject := fmt.Sprintf(TOPIC_PRICES_FORMAT, topic)
	return n.client.Publish(subject, data)
}
