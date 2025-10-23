package main

import (
	"fmt"
	"sync"
	"time"
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
	for v := range intChan {
		fmt.Println(v)
	}

	interfaceChan := make(chan interface{}, 3)
	interfaceChan <- 1
	interfaceChan <- "a"
	p1 := Person{"aaa", 111}
	interfaceChan <- p1
	close(interfaceChan)
	// for v := range interfaceChan {
	// 	fmt.Println(v)
	// }
	<-interfaceChan
	<-interfaceChan
	p2 := <-interfaceChan
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
	//todo关闭之后。读完数据了能不能继续读
	//怎么确认有没有关闭
	// for v := range ch3 {
	// 	fmt.Println(v)
	// }
	for {
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
}

type Person struct {
	Name string
	Age  int
}
