package main

import (
	"demo1/internal/app/mydemo/service"
	"encoding/csv"
	"os"
	"strconv"
	"time"

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

	err = service.LoadConfig()
	if err != nil {
		service.Logger.Error("LoadConfig err", zap.Error(err))

	}

	// 初始化数据库
	err = service.ServiceInitDB(service.Cfg.Database.Dsn)
	if err != nil {
		service.Logger.Error("InitDB err", zap.Error(err))

	}

	timeday := time.Now()
	today := timeday.Year()*10000 + int(timeday.Month())*100 + timeday.Day()
	yestodaytime := time.Now().AddDate(0, 0, -1)
	yestoday := yestodaytime.Year()*10000 + int(yestodaytime.Month())*100 + yestodaytime.Day()

	// uids, err := service.GetCreateUidFromUserCheckinJoin()
	// if err != nil {
	// 	service.Logger.Error("GetCreateUidFromUserCheckinJoin err", zap.Error(err))
	// }

	userCheckinJoins, err := service.GetCreateUidCidFromUserCheckinJoin()
	if err != nil {
		service.Logger.Error("GetCreateUidCidFromUserCheckinJoin err", zap.Error(err))
	}

	service.Logger.Info("uids", zap.Any("uidcid", userCheckinJoins))

	//创建csv
	file, err := os.Create("two_day_not_checkin.csv")
	if err != nil {
		service.Logger.Error("Create err", zap.Error(err))
	}
	defer file.Close()

	//初始化
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入
	writer.Write([]string{
		"uid",
		"cid",
	})

	for _, v := range userCheckinJoins {
		userCheckinRecord1, err := service.GetUserCheckinRecord(v.Uid, v.Cid, today)
		if err != nil {
			service.Logger.Error("GetUserCheckinRecord1 err", zap.Error(err))
		}
		userCheckinRecord2, err := service.GetUserCheckinRecord(v.Uid, v.Cid, yestoday)
		if err != nil {
			service.Logger.Error("GetUserCheckinRecord1 err", zap.Error(err))
		}

		if userCheckinRecord1 == nil && userCheckinRecord2 == nil {
			writer.Write([]string{
				strconv.Itoa(v.Uid),
				strconv.Itoa(v.Cid),
			})
		}
	}

}
