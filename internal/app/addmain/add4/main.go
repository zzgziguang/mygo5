package main

import (
	"demo1/internal/app/mydemo/service"
	"fmt"
	"os"
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

	// timeday := time.Now()
	// today := timeday.Year()*10000 + int(timeday.Month())*100 + timeday.Day()
	yestodaytime := time.Now().AddDate(0, 0, -1)
	yestoday := yestodaytime.Year()*10000 + int(yestodaytime.Month())*100 + yestodaytime.Day()

	// uids, err := service.GetCreateUidFromUserCheckinJoin()
	// if err != nil {
	// 	service.Logger.Error("GetCreateUidFromUserCheckinJoin err", zap.Error(err))
	// }

	page := 1
	pagesize := 3
	isasc := true
	//分页查
	userCheckinJoins, err := service.GetCreateUidCidFromUserCheckinJoin(page, pagesize, isasc)
	if err != nil {
		service.Logger.Error("GetCreateUidCidFromUserCheckinJoin err", zap.Error(err))
	}

	service.Logger.Info("uids", zap.Any("uidcid", userCheckinJoins))

	//创建csv
	//file, err := os.Create("two_day_not_checkin.csv")
	file, err := os.Create("two_day_not_checkin.txt")
	if err != nil {
		service.Logger.Error("Create err", zap.Error(err))
	}
	defer file.Close()

	// //初始化
	// // txt 自己写"\t"分隔 "\n"换行
	// writer := csv.NewWriter(file)
	// defer writer.Flush()
	// // 写入
	// err = writer.Write([]string{
	// 	"uid",
	// 	"cid",
	// })
	// if err != nil {
	// 	service.Logger.Error("Write err", zap.Error(err))
	// }

	num, err := file.WriteString("uid\tcid\n")
	if err != nil {
		service.Logger.Error("WriteString err", zap.Error(err))
	}
	service.Logger.Info("num", zap.Int("num", num))

	for _, v := range userCheckinJoins {
		//date >= <= 限制时间范围
		//select count
		count, err := service.GetUserCheckinRecordTwoDayCount(v.Uid, v.Cid, yestoday)
		if err != nil {
			service.Logger.Error("GetUserCheckinRecord1 err", zap.Error(err))
		}

		if count == 0 {
			str := fmt.Sprintf("%d\t%d\n", v.Uid, v.Cid)
			file.WriteString(str)
		}
	}

}
