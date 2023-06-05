package apm

import (
	"APM-server/internal/pkg/log"
	"APM-server/pkg/kafka"
	"APM-server/pkg/tools"
	"context"
	"github.com/Shopify/sarama"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"net/http"
	"sync"
	"syscall"
)

// var websocketConnections = make(map[int32]*websocket.Conn)
// var websocketMutex sync.Mutex
//var metricsCh = make(chan tools.Metrics, 4)
//var metricsCh chan tools.Metrics
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewConsumePartition(ctx context.Context, wg *sync.WaitGroup, topics []string, partition int32,metricsCh chan tools.Metrics) {
	c, err := kafka.NewConsumer(kafka.KS().Client())
	if err != nil {
		log.Errorw("create consumer error", "err", err)
		return
	}

	for i := 0; i < len(topics); i++ {
		i := i
		p, err := c.ConsumePartition(topics[i], partition, sarama.OffsetNewest)
		if err != nil {
			log.Errorw("create ConsumePartition error", "error", err)
		}
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					//log.Infow("ConsumePartition closing", "topic", topics[i])
					p.Close()
					return
				case msg, ok := <-p.Messages():
					if !ok {
						log.Errorw("error fetch msg from", "topic", topics[i])
						continue
					}
					met, _ := tools.ProcessBytesMessage(msg.Value)
					//log.Infow("Received mes", "partition", msg.Partition, "met:", met, "offset", msg.Offset)
					metricsCh <- met
				case err := <-p.Errors():
					log.Errorw("error receive msg", "err", err)
				}
			}
		}()
	}
}
func WebsocketHandler(c *gin.Context) {
	log.Infow("starting a new websocket connection","X-Request-ID",c.GetString("X-Request-ID"))
	metricsCh := make(chan tools.Metrics, 4)
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Errorw("升级为 WebSocket 连接失败:", "err", err)
		return
	}

	// 获取分区信息
	addrStr := c.Param("addr")
	addrPartitionMap := viper.GetStringMap("kafka.map")
	originPartition := addrPartitionMap[addrStr]
	var partition int32
	switch v := originPartition.(type) {
	case int:
		partition = int32(v)
	case int32:
		partition = v
	case int64:
		partition = int32(v)
	// 添加其他可能的类型转换
	default:
		log.Errorw("无法将addr映射的partition转换为int32:", "v", v)
		return
	}

	topics := viper.GetStringSlice("kafka.topics")
	var wg sync.WaitGroup
	wg.Add(len(topics))
	ctx, cancel := context.WithCancel(context.Background())
	go NewConsumePartition(ctx, &wg, topics, partition,metricsCh)
	go func() {
		for {
			select {
			case met := <-metricsCh:
				// 在写入数据之前先进行一次写入测试
				err := conn.WriteMessage(websocket.PingMessage, []byte{})
				if err != nil {
					if errors.Is(err,syscall.EPIPE) {
						log.Infow("websocket connection closing","X-Request-ID",c.GetString("X-Request-ID"))
					}else {
						log.Errorw("write ping error,connection closing", "err", err,"X-Request-ID",c.GetString("X-Request-ID"))
					}
					cancel()
					wg.Wait()
					conn.Close()
					//close(metricsCh)
					log.Infow("websocket connection successfully closed","X-Request-ID",c.GetString("X-Request-ID"))
					return
				}
				err = conn.WriteJSON(met)
				if err != nil {
					log.Errorw("send msg error", "err", err)
					continue
				}
			}
		}
	}()

}
