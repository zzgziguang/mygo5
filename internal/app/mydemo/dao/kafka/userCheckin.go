package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
)

// 添加kafka的数据
func ProducerSend(uid int) (partition int32, offset int64, err error) {
	struid := fmt.Sprintf("恭喜%d打卡成功", uid)
	msg := &sarama.ProducerMessage{
		Topic: "topic_user_checkin",
		Value: sarama.StringEncoder(struid),
	}
	partition, offset, err = KafkaClient.SendMessage(msg)
	return
}

//消费者数据
