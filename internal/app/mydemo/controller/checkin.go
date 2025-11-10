package controller

import (
	"context"
	"demo1/internal/app/mydemo/model"
	"demo1/internal/app/mydemo/service"
	"sort"
	"strconv"
	"sync"

	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AddCheckinHandler(c *gin.Context) {
	title := c.PostForm("title")
	if title == "" {
		MakeApiResponse(c, 1001, "标题不能为空")
		return
	}

	strweight := c.PostForm("weight")
	if strweight == "" {
		MakeApiResponse(c, 1001, "weight不能为空")
		return
	}

	weight, err := strconv.Atoi(strweight)
	if err != nil {
		service.Logger.Error("Atoistrweight err", zap.Error(err))
		MakeApiResponse(c, 1001, "weight类型转换错误"+err.Error())
		return
	}

	createat := time.Now()
	checkinstatus := model.CheckinNormal

	newCheckin := &model.Checkin{
		Title:         title,
		CreateAt:      &createat,
		UpdateAt:      &createat,
		CheckinStatus: checkinstatus,
	}

	//插入数据库
	result := service.CreateCheckin(newCheckin)
	if result.Error != nil {
		service.Logger.Error("CreateCheckin err", zap.Error(result.Error))
		MakeApiResponse(c, 1, "插入数据库错误"+result.Error.Error())
		return
	}

	//更新参与打卡的权重
	err = service.UpdateCheckinWeight(newCheckin, weight)
	if err != nil {
		service.Logger.Error("UpdateCheckinWeight err", zap.Error(err))
		MakeApiResponse(c, 1, "更新参与人数失败"+err.Error())
		return
	}

	// 返回成功响应
	MakeApiResponse(c, 0, newCheckin)
}

// 查询所有ckeckin根据id排序分页
func GetCheckinHandlerAll(c *gin.Context) {

	strpage := c.Query("page")
	order := c.Query("order")
	struid := c.Query("uid")

	if strpage == "" {
		MakeApiResponse(c, 1001, "page不能为空")
		return
	}

	if order != "desc" && order != "asc" {
		MakeApiResponse(c, 1001, "order只能是desc或asc")
		return
	}

	if struid == "" {
		MakeApiResponse(c, 1001, "uid不能为空")
		return
	}

	page, err := strconv.Atoi(strpage)
	if err != nil {
		service.Logger.Error("pageAtoi err", zap.Error(err))
		MakeApiResponse(c, 1001, "page必须为数字"+err.Error())
		return
	}
	if page <= 0 {
		service.Logger.Error("page err", zap.Error(err))
		MakeApiResponse(c, 1001, "page必须为正整数")
		return
	}

	uid, err := strconv.Atoi(struid)
	if err != nil {
		service.Logger.Error("uidAtoi err", zap.Error(err))
		MakeApiResponse(c, 1001, "uid类型转换错误"+err.Error())
		return
	}

	timeday := time.Now()
	date := timeday.Year()*10000 + int(timeday.Month())*100 + timeday.Day()

	var wg sync.WaitGroup
	var mu sync.Mutex

	var checkinSlice []model.Checkin
	var joinSlice []model.UserCheckinJoin
	var recordSlice []model.UserCheckinRecord

	var rank int
	var zrankm map[string]int
	var err1 error
	var err2 error
	var err3 error
	var err4 error

	var cidNumMap map[int]int = make(map[int]int, 0)
	var cidNum int

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//获取全部
	//checkinData，checkinList,checkinSlice,checkinResult,checkinRes
	wg.Add(1)
	go func() {
		defer wg.Done()

		defer func() {
			err := recover()
			if err != nil {
				service.Logger.Error("GetCheckinAll panic", zap.Any("panic", err))
			}
		}()

		checkinSlice, err1 = service.GetCheckinAll()
		if err1 != nil {
			service.Logger.Error("err1", zap.Error(err1))
			return
		}

		mu.Lock()
		for _, v := range checkinSlice {
			cidNumMap[v.Id]++ //加锁
		}
		mu.Unlock()
	}()

	//获取打卡名次zset缓存
	wg.Add(1)
	go func() {
		defer wg.Done()

		defer func() {
			err := recover()
			if err != nil {
				service.Logger.Error("GetZsetCheckinNum panic", zap.Any("panic", err))
				cancel()
			}
		}()

		zrankm, err2 = service.GetZsetCheckinNum()
		if err2 != nil {
			service.Logger.Error("err2", zap.Error(err2))
			cancel()
		}

	}()

	//根据uid查询join表
	wg.Add(1)
	go func() {
		defer wg.Done()

		defer func() {
			err := recover()
			if err != nil {
				service.Logger.Error("GetUserCheckinJoinByuid panic", zap.Any("panic", err))
			}
		}()

		joinSlice, err3 = service.GetUserCheckinJoinByuid(uid)
		if err3 != nil {
			service.Logger.Error("err3", zap.Error(err3))
			return
		}

		mu.Lock()
		for _, v := range joinSlice { //并发写会有问题,应该加锁，这里为什么没有报错
			cidNumMap[v.Cid]++ //加锁
		}
		mu.Unlock()
		// i := 0
		// fmt.Print(uid / i)
	}()

	var weatherErr error
	var todayWeather service.Weather
	wg.Add(1)
	go func() {
		defer wg.Done()

		defer func() {
			err := recover()
			if err != nil {
				service.Logger.Error("GetWeather panic", zap.Any("panic", err))
			}
		}()

		weatherCtx, weatherCancel := context.WithTimeout(ctx, 200*time.Millisecond) //50*time.Millisecond
		defer weatherCancel()

		todayWeather, weatherErr = service.GetWeather(weatherCtx, "北京")
		if weatherErr != nil {
			service.Logger.Error("weatherErr", zap.Error(weatherErr))
			return
		}

		service.Logger.Info("today weather", zap.Any("todayWeather", todayWeather))
	}()

	wg.Wait()

	if err1 != nil {
		service.Logger.Error("err2", zap.Error(err2))
		MakeApiResponse(c, 1, "redis查询失败"+err1.Error())
		return
	}

	if err2 != nil {
		service.Logger.Error("err2", zap.Error(err2))
		MakeApiResponse(c, 1, "redis查询失败"+err2.Error())
		return
	}

	//根据uid获取用户所有的已参与打卡，转换为map[cid]Join，和今日所有的打卡记录转换为map[cid]Record

	if err3 != nil {
		service.Logger.Error("err3", zap.Error(err3))
		MakeApiResponse(c, 1, "查询参与数据库的错误"+err3.Error())
		return
	}

	if weatherErr != nil {
		MakeApiResponse(c, 2002, "超时"+weatherErr.Error())
		service.Logger.Error("weatherErr", zap.Error(weatherErr))
		return
	}

	//获取用户今日打卡记录
	//recordSlice, err4 = service.GetUserCheckinRecordByUidDate(uid, date)
	var cidSlice []int
	for _, v := range checkinSlice {
		cidSlice = append(cidSlice, v.Id)
	}

	recordSlice, err4 = service.GetUserCheckinRecordInCheckinId(uid, date, cidSlice)
	// 使用in语法 cid in join的cid
	if err4 != nil {
		service.Logger.Error("err4", zap.Error(err4))
		MakeApiResponse(c, 1, "查询打卡记录数据库的错误"+err4.Error())
		return
	}

	userCheckinJoinMap := make(map[int]*model.UserCheckinJoin, 0)
	for _, userCheckinJoin := range joinSlice {
		userCheckinJoinMap[userCheckinJoin.Cid] = &userCheckinJoin
	}

	userCheckinRecordMap := make(map[int]*model.UserCheckinRecord, 0)
	for _, userCheckinRecord := range recordSlice {
		userCheckinRecordMap[userCheckinRecord.Cid] = &userCheckinRecord
	}

	responsecheckin := make([]model.ResponseCheckinItem, 0)

	//var responsecheckin []model.ResponseCheckinItem=make([]model.ResponseCheckinItem, 3,3)
	//responsecheckin := make([]model.ResponseCheckinItem,0)
	var checkinSortList model.CheckinSlice
	for _, v := range checkinSlice {
		cid := v.Id
		joinBool := false
		recordBool := false
		joinTime := time.Time{}

		userCheckinJoin, ok := userCheckinJoinMap[cid]
		if ok {
			joinBool = true
			if userCheckinJoin.CreateAt != nil {
				joinTime = *userCheckinJoin.CreateAt
			}
		}

		_, ok = userCheckinRecordMap[cid]
		if ok {
			recordBool = true
		}

		rankm, ok := zrankm[strconv.Itoa(cid)]
		if ok {
			rank = rankm
		}

		checkinSortList = append(checkinSortList, model.CheckinSort{
			Checkin:    v,
			JoinBool:   joinBool,
			RecordBool: recordBool,
			JoinTime:   joinTime,
			Rank:       rank,
		})
	}

	// 排序
	sort.Sort(checkinSortList)

	for _, v := range checkinSortList {
		//var weather map[string]interface{}
		cid := v.Checkin.Id
		cidNumm, ok := cidNumMap[cid]
		if ok {
			cidNum = cidNumm
		}
		// if v.RecordBool == true {

		// }
		checkinre := model.ResponseCheckinItem{
			Id:            v.Checkin.Id,
			Title:         v.Checkin.Title,
			CreateAt:      v.Checkin.CreateAt.Format("2006年01月02日 15点04分05秒"),
			UpdateAt:      v.Checkin.UpdateAt.Format("2006年01月02日 15点04分05秒"),
			CheckinStatus: v.Checkin.CheckinStatus,
			JoinBool:      v.JoinBool,               //是否参与
			RecordBool:    v.RecordBool,             //今日是否打卡
			JoinNumber:    int64(v.Checkin.JoinNum), //参与人数
			Rank:          v.Rank,                   //参与人数名次
			Weight:        v.Checkin.Weight,         //打卡的权重
			JoinTime:      v.JoinTime.Format("2006年01月02日 15点04分05秒"),
			CidNum:        cidNum,
		}

		responsecheckin = append(responsecheckin, checkinre)
	}

	datastruct := service.DataStruct{
		Weather: model.WeatherItem{
			Temperature: todayWeather.Temperature,
			Weather:     todayWeather.Weather,
		},
		Slices: responsecheckin,
	}

	MakeApiResponse(c, 0, datastruct)
}
