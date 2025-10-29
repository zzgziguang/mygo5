package main

import (
	"context"
	"demo1/internal/app/mydemo/service"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

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
	for {
		// 启动消费循环
		err = consumerGroup.Consume(ctx, []string{topic}, &consumerGroupHandler{})
		if err != nil {
			service.Logger.Error("Consume err", zap.Error(err))
			break
		}
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-sigChan
	service.Logger.Info("sigChan:", zap.String("sigChan", fmt.Sprintf("获得信号%v", sig)))

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
