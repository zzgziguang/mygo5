package main

import (
	"fmt"
	"time"
)

func main() {

	int1 := make(chan int, 1)
	int1 <- 10
	for i := 1; i <= 10; i++ {

		select {
		case i2 := <-int1:
			fmt.Println("i2", i2)
		case i1 := <-int1:
			fmt.Println(i1)
		case <-time.After(3 * time.Second):
			fmt.Println("timeafter")

		default:
			fmt.Println("default")
		}

	}
}
