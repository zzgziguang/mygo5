package main

import (
	"bufio"
	"demo1/internal/app/mydemo/service"
	"os"
	"strings"

	"go.uber.org/zap"
)

func main() {
	var err error

	err = service.LoggerInit()
	if err != nil {
		panic(err)
	}
	//Logger, err = zap.NewDevelopment()
	defer service.SyncLogger()

	// 打开这个 csv 文件
	file, err := os.Open("../add4/two_day_not_checkin.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	service.Logger.Info("uidcid", zap.String("uidcid", "逐行读取"))

	str := scanner.Text()
	service.Logger.Info("uidcid", zap.String("uidcid", str))

	err = scanner.Err()
	if err != nil {
		service.Logger.Error("NewScanner err", zap.Error(err))
		panic(err)
	}

	for {
		bol := scanner.Scan()
		if bol == true {
			str := scanner.Text()
			strs := strings.Split(str, "\t")
			service.Logger.Info("uidcid", zap.String("uid", strs[0]), zap.String("cid", strs[1]))
		} else {
			service.Logger.Error("Scan err", zap.String("err", "读取完"))
			break
		}
	}

	// // 初始化
	// reader := csv.NewReader(file)
	// //返回数据
	// //不用readall
	// //uidcid
	// record, err := reader.Read()
	// if err != nil {
	// 	service.Logger.Error("Read err", zap.Error(err))
	// 	panic(err)
	// }
	// service.Logger.Info("uidcid", zap.String("uid", record[0]), zap.String("cid", record[1]))

	// for {
	// 	record, err := reader.Read()
	// 	if err != nil {
	// 		service.Logger.Error("Read err", zap.Error(err))
	// 		return
	// 	}
	// 	service.Logger.Info("uidcid", zap.String("uid", record[0]), zap.String("cid", record[1]))
	// }

}
