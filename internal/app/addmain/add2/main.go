package main

import (
	"context"
	"fmt"
	"time"
)

var NameKey string = "name"

// func worker(ctx context.Context, cancel context.CancelFunc) {
func worker(ctx context.Context) {
	name := ctx.Value(NameKey)
	if name != nil {
		fmt.Println(name)
	}

	for i := 1; i <= 3; i++ {
		if i == 2 {
			time.Sleep(300 * time.Millisecond)
			//cancel()
		}
		select {
		case <-time.After(300 * time.Millisecond):
			fmt.Println("worker", i, time.Now())
		case <-ctx.Done():
			fmt.Println("cancel", ctx.Err(), time.Now())
		}
	}
	fmt.Println("nihao")
}

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, NameKey, "lisi")

	// withCancelCtx, withCancelCancel := context.WithCancel(ctx)
	// defer withCancelCancel()
	// go worker(withCancelCtx, withCancelCancel)

	withTimeOutCtx, withTimeOutCancel := context.WithTimeout(ctx, 1*time.Second)
	defer withTimeOutCancel()
	go worker(withTimeOutCtx)

	// withDeadlineCtx, withDeadlineCancel := context.WithDeadline(ctx, time.Now().Add(1*time.Second))
	// defer withDeadlineCancel()
	// go worker(withDeadlineCtx)

	time.Sleep(3 * time.Second)

}
