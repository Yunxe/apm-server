package kafka

import (
	"APM-server/internal/pkg/log"
	"github.com/Shopify/sarama"
	"github.com/spf13/viper"
	"strconv"
	"time"
)

func InitTopics() {
	topics := viper.GetStringSlice("kafka.topics")
	numPartitions := viper.GetInt32("kafka.number-partition")

	t, _ := KS().Client().Topics()
	log.Infow("existed topics:", "topics", t)
	for _, v := range t {
		partitions, _ := KS().Client().Partitions(v)
		log.Infow("existed partitions from topic", "topic", v, "partitions", partitions)
	}

	admin, err := sarama.NewClusterAdminFromClient(KS().Client())
	if err != nil {
		log.Errorw("new cluster admin error", "err", err)
	}
	defer admin.Close()

	topicDetail := &sarama.TopicDetail{
		ReplicationFactor: -1,
		NumPartitions:     numPartitions,
		ConfigEntries: map[string]*string{
			"retention.ms": func(s string) *string { return &s }(strconv.FormatInt(int64(time.Minute*3), 10)),
		},
	}
	if len(t) < int(numPartitions) {
		for _, t := range topics {
			err := admin.CreateTopic(t, topicDetail, false)
			if err != nil {
				log.Errorw("create topic error", "err", err)
			} else {
				log.Infow("successfully create topic", "topic", t)
			}
		}
	}

}
