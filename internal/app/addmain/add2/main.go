package main

import (
	"demo1/internal/app/mydemo/controller"
	"demo1/internal/app/mydemo/service"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
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

	err = service.LoadConfig()
	if err != nil {
		service.Logger.Error("LoadConfig err", zap.Error(err))

	}

	// 初始化数据库
	err = service.ServiceInitDB(service.Cfg.Database.Dsn)
	if err != nil {
		service.Logger.Error("InitDB err", zap.Error(err))

	}

	//初始化 Redis
	service.ServiceInitRedis(service.Cfg.Redis.Addr, service.Cfg.Redis.Password, service.Cfg.Redis.DB)

	//初始化 kafka
	err = service.ServiceInitKafka()
	if err != nil {
		service.Logger.Error("InitKafka err", zap.Error(err))
		panic(err)
	}
	defer service.Closekafka()

	// 1. 启动服务到后台
	go func() {
		r := gin.Default()
		r.POST("/api/userCheckin/add", controller.AddUserCheckinHandler)
		r.Run(":8081")
	}()

	time.Sleep(300 * time.Millisecond)

	var rankMapChan chan map[string]interface{}
	var rankMap map[string]interface{}
	var wg sync.WaitGroup
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
			defer resp.Body.Close()

			checkin, err := io.ReadAll(resp.Body)
			if err != nil {
				service.Logger.Error("ReadAll err", zap.Error(err))
				return
			}

			var apiJson ApiJson
			err = json.Unmarshal(checkin, &apiJson)
			if err != nil {
				service.Logger.Error("Unmarshal err", zap.Error(err))
				return
			}

			rank = apiJson.Data.Rank
			rankMap = map[string]interface{}{
				"uid":  strconv.Itoa(15 + i),
				"rank": rank, //今日打卡名次
			}
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
