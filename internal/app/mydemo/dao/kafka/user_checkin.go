package kafka

import (
	"demo1/internal/app/mydemo/model"
	"encoding/json"

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

// 添加usercheckinjoin kafka的数据
func ProduceKafkaUserCheckinJoinMessage(userCheckinJoinMsg *model.UserCheckinJoin) (partition int32, offset int64, err error) {
	value, err := json.Marshal(userCheckinJoinMsg)
	if err != nil {
		return
	}

	msg := &sarama.ProducerMessage{
		Topic: "topic_user_checkin",
		Value: sarama.StringEncoder(value),
	}

	partition, offset, err = KafkaClient.SendMessage(msg)
	return
}

// 添加usercheckinrecord kafka的数据
func ProduceKafkaUserCheckinRecordMessage(userCheckinRecordMsg *model.UserCheckinRecord) (partition int32, offset int64, err error) {
	value, err := json.Marshal(userCheckinRecordMsg)
	if err != nil {
		return
	}

	msg := &sarama.ProducerMessage{
		Topic: "topic_user_checkin",
		Value: sarama.StringEncoder(value),
	}

	partition, offset, err = KafkaClient.SendMessage(msg)
	return
}
