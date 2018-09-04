package rabbit

type Configuration struct {
	AmqpUrl           string
	Exchange          string
	ExchangeType      string
	ConnectRetryLimit int
}

func NewConfiguration(amqpUrl string, exchange string, exchangetype string, retry int) Configuration {
	return Configuration{
		AmqpUrl:           amqpUrl,
		Exchange:          exchange,
		ExchangeType:      exchangetype,
		ConnectRetryLimit: retry,
	}
}
