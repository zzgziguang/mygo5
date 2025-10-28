package main

import (
	"context"
	"demo1/internal/app/mydemo/service"
	"sync"

	"github.com/IBM/sarama"
	"go.uber.org/zap"
)

var valueChan = make(chan string, 10)

func main() {
	var err error

	err = service.LoggerInit()
	if err != nil {
		panic(err)
	}
	//Logger, err = zap.NewDevelopment()
	defer service.SyncLogger()

	brokers := []string{"localhost:9092"}
	topic := "topic_user_checkin"
	groupID := "group_checkin"

	config := sarama.NewConfig()
	// config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		service.Logger.Error("err", zap.Error(err))
	}
	defer func() {
		err := consumerGroup.Close()
		if err != nil {
			service.Logger.Error("err", zap.Error(err))
		}
	}()

	var wg sync.WaitGroup
	// for i := 1; i <= 3; i++ {
	// }
	//TODO 写个每秒一条数据的生产者
	//TODO 整理chan close的情况
	//TODO 生产者uid，cid，时间戳 json
	//TODO go signal 使用用法
	wg.Add(1)
	go func() {
		defer wg.Done()
		for msg := range valueChan {
			service.Logger.Info("v1", zap.String("v value", msg))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for msg := range valueChan {
			service.Logger.Info("v2", zap.String("v value", msg))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for msg := range valueChan {
			service.Logger.Info("v3", zap.String("v value", msg))
		}
	}()
	ctx := context.Background()

	// 启动消费循环
	err = consumerGroup.Consume(ctx, []string{topic}, &consumerGroupHandler{})
	if err != nil {
		service.Logger.Error("Consume err", zap.Error(err))
	}
	close(valueChan)

	wg.Wait()

}

type consumerGroupHandler struct{}

// Setup 在每个会话开始前调用
func (h consumerGroupHandler) Setup(session sarama.ConsumerGroupSession) (err error) {
	return nil
}

// Cleanup 在每个会话结束后调用
func (h consumerGroupHandler) Cleanup(claim sarama.ConsumerGroupSession) (err error) {
	return nil
}

// ConsumeClaim 处理分配给该消费者的分区中的消息
func (h consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) (err error) {
	for message := range claim.Messages() {
		service.Logger.Info("message", zap.String("topic", message.Topic), zap.Int32("partition", message.Partition),
			zap.Int64("offset", message.Offset), zap.String("value", string(message.Value)))

		valueChan <- string(message.Value)
		session.MarkMessage(message, "")
	}
	return nil
}
