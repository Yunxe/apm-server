package kafka

import (
	"github.com/Shopify/sarama"
)



type KafkaOptions struct {
	ConsumerReturnErr bool
	GroupID           string
	Brokers           []string
}


func NewKafkaClient(opts *KafkaOptions) (*sarama.Client, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = opts.ConsumerReturnErr
	client, err := sarama.NewClient(opts.Brokers, config)
	if err != nil {
		return nil, err
	}
	return &client, nil
}


func StoreClient(client *sarama.Client) {
	once.Do(func() {
		kafkaStore.client = *client
	})
}




