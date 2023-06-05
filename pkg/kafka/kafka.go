package kafka

import (
	"github.com/Shopify/sarama"
	"sync"
)
var (
	once       sync.Once
	kafkaStore = &KafkaStore{}
)

type KafkaStore struct {
	client sarama.Client
	//cg     sarama.ConsumerGroup
}

func KS() *KafkaStore {
	return kafkaStore
}

//func (ks KafkaStore) CG() sarama.ConsumerGroup {
//	return ks.cg
//}

func (ks KafkaStore) Client() sarama.Client {
	return ks.client
}
