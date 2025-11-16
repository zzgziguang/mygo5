package main

import (
	"demo1/internal/app/mydemo/service"
	"encoding/csv"
	"os"

	"go.uber.org/zap"
)

func main() {
	var err error

	err = service.LoggerInit()
	if err != nil {
		//panic(err)
	}
	//Logger, err = zap.NewDevelopment()
	defer service.SyncLogger()
	// 打开这个 csv 文件
	file, err := os.Open("../add4/two_day_not_checkin.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 初始化
	reader := csv.NewReader(file)
	//返回数据
	records, err := reader.ReadAll()
	if err != nil {
		service.Logger.Error("ReadAll err", zap.Error(err))
		panic(err)
	}

	if len(records) <= 1 {
		service.Logger.Error("records err", zap.String("err", "csn为空"))
		return
	}

	for i, ids := range records {
		if i == 0 {
			continue
		}
		service.Logger.Info("uids", zap.String("uid", ids[0]), zap.String("cid", ids[1]))
	}

}
