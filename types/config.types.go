package types

type Config struct {

	// Rabbit
	RabbitHost        string
	RabbitPort        int
	RabbitUser        string
	RabbitPass        string
	RabbitVHost       string
	RabbitExchange    string
	RabbitExchangeKey string
	RabbitQueue       string
}
