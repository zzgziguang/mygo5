package kafka

import (
	"github.com/IBM/sarama"
)

// 添加kafka的数据
func ProducerSend(value []byte) (partition int32, offset int64, err error) {
	msg := &sarama.ProducerMessage{
		Topic: "topic_user_checkin",
		Value: sarama.StringEncoder(value),
	}

	partition, offset, err = KafkaClient.SendMessage(msg)
	return
}
