package main

import (
	"demo1/internal/app/mydemo/service"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"go.uber.org/zap"
)

var NameKey string = "name"

// func worker(ctx context.Context, cancel context.CancelFunc) {
// func worker(ctx context.Context) {
// 	name := ctx.Value(NameKey)
// 	if name != nil {
// 		fmt.Println(name)
// 	}

// 	for i := 1; i <= 3; i++ {
// 		if i == 2 {
// 			time.Sleep(300 * time.Millisecond)
// 			//cancel()
// 		}
// 		select {
// 		case <-time.After(300 * time.Millisecond):
// 			fmt.Println("worker", i, time.Now())
// 		case <-ctx.Done():
// 			fmt.Println("cancel", ctx.Err(), time.Now())
// 		}
// 	}
// 	fmt.Println("nihao")
// }

type ApiJson struct {
	Code    int    `json:"code"`
	Data    Data   `json:"data"`
	Message string `json:"message"`
}

type Data struct {
	Rank int64 `json:"rank"`
}

func main() {
	var err error
	var rank int64

	err = service.LoggerInit()
	if err != nil {
		panic(err)
	}
	//Logger, err = zap.NewDevelopment()
	defer service.SyncLogger()

	rankMapChan := make(chan map[string]interface{}, 10)
	var rankMap map[string]interface{}
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			urlkey := "http://localhost:8081/api/userCheckin/add"
			data := url.Values{}
			data.Set("uid", strconv.Itoa(15+i))
			data.Set("cid", "1")

			strEncode := data.Encode()
			resp, err := http.Post(urlkey, "application/x-www-form-urlencoded", strings.NewReader(strEncode))
			if err != nil {
				service.Logger.Error("Post err", zap.Error(err))
				return
			}
			//service.Logger.Info("resp", zap.Any("resp", resp.Body))
			defer resp.Body.Close()

			checkin, err := io.ReadAll(resp.Body)
			if err != nil {
				service.Logger.Error("ReadAll err", zap.Error(err))
				return
			}
			//service.Logger.Info("checkin", zap.Any("checkin", checkin))

			var apiJson ApiJson
			err = json.Unmarshal(checkin, &apiJson)
			if err != nil {
				service.Logger.Error("Unmarshal err", zap.Error(err))
				return
			}
			//service.Logger.Info("apijson", zap.Any("apijson", apiJson))

			code := apiJson.Code
			if code != 0 {
				service.Logger.Error("code err", zap.Int("code", code), zap.String("msg", apiJson.Message))
				return
			}
			mu.Lock()
			rank = apiJson.Data.Rank

			rankMap = map[string]interface{}{
				"uid":  strconv.Itoa(15 + i),
				"rank": rank, //今日打卡名次
			}
			mu.Unlock()

			rankMapChan <- rankMap
		}(i)
	}
	wg.Wait()
	close(rankMapChan)
	for rankMap := range rankMapChan {
		service.Logger.Info("uidCheckinRank", zap.Any("rankmap", rankMap))
	}

}

// return
// ctx := context.Background()
// ctx = context.WithValue(ctx, NameKey, "lisi")

// // withCancelCtx, withCancelCancel := context.WithCancel(ctx)
// // defer withCancelCancel()
// // go worker(withCancelCtx, withCancelCancel)

// withTimeOutCtx, withTimeOutCancel := context.WithTimeout(ctx, 1*time.Second)
// defer withTimeOutCancel()
// go worker(withTimeOutCtx)

// // withDeadlineCtx, withDeadlineCancel := context.WithDeadline(ctx, time.Now().Add(1*time.Second))
// // defer withDeadlineCancel()
// // go worker(withDeadlineCtx)

// time.Sleep(3 * time.Second)

//}
