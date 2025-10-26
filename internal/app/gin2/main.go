package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	var intChan chan int
	intChan = make(chan int, 2)
	fmt.Printf("%v,%p\n", intChan, &intChan)

	intChan <- 10
	num := 5
	intChan <- num
	// intChan <- 4 //与容量有关，多了不行
	fmt.Printf("%v,%p\n", intChan, &intChan)
	close(intChan) //guanbi

	// n1 := <-intChan
	// n2, ok := <-intChan
	// fmt.Println(n1, n2, ok)
	//取了不能再取
	// for v := range intChan {
	// 	fmt.Println(v)
	// }
	for i := 0; i < 5; i++ {
		select {
		case v, ok := <-intChan:
			if ok {
				fmt.Println("intChan数据:", v)
			} else {
				fmt.Println("intChan已关闭")
			}
		default:
			fmt.Println("intChan无数据")
		}
	}

	interfaceChan := make(chan interface{}, 3)
	interfaceChan <- 1
	interfaceChan <- "a"
	p1 := Person{"aaa", 111}
	interfaceChan <- p1
	//interfaceChan <- Person{"bbb", 222}
	close(interfaceChan)
	// for v := range interfaceChan {
	// 	fmt.Println(v)
	// }
	<-interfaceChan
	<-interfaceChan
	p2 := <-interfaceChan
	//p3 := <-interfaceChan
	fmt.Println(p2)
	ch1 := make(chan int, 1)
	ch1 <- 1
	ch2 := make(chan string, 1)
	ch2 <- "aaa"
	select {
	case msg := <-ch1:
		fmt.Println(msg)
	default:
		fmt.Println(ch2)
	}

	//goroutine
	var wg sync.WaitGroup
	wg.Add(1)
	ch3 := make(chan string, 2)
	go func() {
		defer wg.Done()
		str1 := "aaa"
		ch3 <- str1
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		str2 := "bbb"
		time.Sleep(3 * time.Second)
		ch3 <- str2
	}()
	wg.Wait()
	//close(ch3)
	//关闭之后。读完数据了能不能继续读
	//怎么确认有没有关闭
	// for v := range ch3 {
	// 	fmt.Println(v)
	// }
	for i := 0; i < 5; i++ {
		select {
		case msg1 := <-ch3:
			fmt.Println("aaa", msg1)
		case <-time.After(1 * time.Second):
			fmt.Println("超时")
		}
		fmt.Println(time.Now())
	}

	// str3 := <-ch3
	// fmt.Println("str3", str3)

	//超时
	ch4 := make(chan int)
	go func() {
		time.Sleep(3 * time.Second)
		ch4 <- 111
	}()
	select {
	case msg1 := <-ch4:
		fmt.Println(msg1)
	case <-time.After(1 * time.Second):
		fmt.Println("超时")
	}

	const broker = "localhost:9092"
	const topic = "simplified-topic"
	const group = "simplified-group"

	// 启动消费者（goroutine）
	go func() {
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{broker},
			Topic:   topic,
			GroupID: group, // 启用消费者组
		})
		defer reader.Close()

		fmt.Println("消费者已启动，等待消息...")

		for {
			msg, err := reader.ReadMessage(context.Background())
			if err != nil {
				fmt.Println("消费错误:", err)
				continue
			}
			fmt.Printf("收到: %s (分区=%d, offset=%d)\n",
				string(msg.Value), msg.Partition, msg.Offset)

			// 提交 offset
			reader.CommitMessages(context.Background(), msg)
		}
	}()

	// 给消费者一点时间启动
	time.Sleep(500 * time.Millisecond)

	//生产者发送消息
	writer := &kafka.Writer{
		Addr:  kafka.TCP(broker),
		Topic: topic,
	}
	defer writer.Close()

	fmt.Println("生产者发送 3 条消息...")
	for i := 1; i <= 3; i++ {
		msg := kafka.Message{
			Value: []byte(fmt.Sprintf("消息 %d", i)),
		}
		writer.WriteMessages(context.Background(), msg)
		fmt.Printf("发送: 消息 %d\n", i)
		time.Sleep(300 * time.Millisecond)
	}

	// 等待消费完成
	time.Sleep(2 * time.Second)
	fmt.Println("演示结束")
}

type Person struct {
	Name string
	Age  int
}
