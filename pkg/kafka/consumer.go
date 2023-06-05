package kafka

import (
	"github.com/Shopify/sarama"
)

func NewConsumer(client sarama.Client) (sarama.Consumer, error) {
	c,err:=sarama.NewConsumerFromClient(client)
	if err != nil {
		return nil, err
	}
	return c, nil
}

//func StoreConsumerGroup(cg *sarama.ConsumerGroup) {
//	once.Do(func() {
//		kafkaStore.cg = *cg
//	})
//}
