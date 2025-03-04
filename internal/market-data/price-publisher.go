package marketdata

const (
	TOPIC_PRICES_FORMAT = "prices.%s"
)

type PricePublisher interface {
	Publish(topic string, data []byte) error
}
