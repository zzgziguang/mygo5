package kafka

import (
	"github.com/IBM/sarama"
)

var KafkaClient sarama.SyncProducer

// 连接kafka
func DaoInitKafka() (err error) {
	brokerlist := []string{"localhost:9092"}
	config := sarama.NewConfig()
	// config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	// config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true // 成功交付的消息将在success channel返回

	KafkaClient, err = sarama.NewSyncProducer(brokerlist, config)
	if err != nil {
		return
	}
	return
}

// 关闭
func Closekafka() {
	KafkaClient.Close()
}
