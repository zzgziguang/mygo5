package controller

import (
	"context"
	"demo1/internal/app/mydemo/model"
	"demo1/internal/app/mydemo/service"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"sync"

	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type weatherstruct struct {
	Reason     string
	Error_code int
	Result     weatherresult
}

type weatherresult struct {
	City     string
	Realtime realtimeweather
	Future   []futureweather
}

type realtimeweather struct {
	Temperature string
	Humidity    string
	Info        string
	Wid         string
	Direct      string
	Power       string
	Aqi         string
}

type futureweather struct {
	Date        string
	Temperature string
	Weather     string
	Wid         widweather
	Direct      string
}

type widweather struct {
	Day   string
	Night string
}

type datastruct struct {
	Weather model.WeatherItem
	Slices  []model.ResponseCheckinItem
}

type WeatherResult struct {
	Temperature string
	Weather     string
}

func AddCheckinHandler(c *gin.Context) {
	title := c.PostForm("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Error:   "标题不能为空",
		})
		return
	}

	strweight := c.PostForm("weight")
	if strweight == "" {
		c.JSON(http.StatusNotFound, model.APIResponse{
			Success: false,
			Error:   "weight不能为空",
		})
		return
	}

	weight, err := strconv.Atoi(strweight)
	if err != nil {
		c.JSON(http.StatusNotFound, model.APIResponse{
			Success: false,
			Error:   "weight类型转换错误",
		})
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
		service.Logger.Error("CreateCheckin错误", zap.Any("newCheckin", newCheckin), zap.String("err1", result.Error.Error()))
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Success: false,
			Error:   "插入数据库错误",
		})
		return
	}

	//更新参与打卡的权重
	err = service.UpdateCheckinWeight(newCheckin, weight)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Success: false,
			Error:   "更新参与人数失败",
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Message: "打卡添加成功",
		Data:    newCheckin,
	})
}

// 查询所有ckeckin根据id排序分页
func GetCheckinHandlerAll(c *gin.Context) {

	strpage := c.Query("page")
	order := c.Query("order")
	struid := c.Query("uid")

	if strpage == "" {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Error:   "page不能为空",
		})
		return
	}

	if order != "desc" && order != "asc" {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Error:   "order只能是desc或asc",
		})
		return
	}

	if struid == "" {
		c.JSON(http.StatusNotFound, model.APIResponse{
			Success: false,
			Error:   "uid不能为空",
		})
		return
	}

	page, err := strconv.Atoi(strpage)
	if err != nil {
		service.Logger.Error("page格式错误", zap.String("err:", err.Error()))
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Error:   "page必须为数字" + err.Error(),
		})
		return
	}
	if page <= 0 {
		service.Logger.Error("page错误", zap.String("err:", "page值错误"))
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Error:   "page必须为正整数",
		})
		return
	}

	uid, err := strconv.Atoi(struid)
	if err != nil {
		c.JSON(http.StatusNotFound, model.APIResponse{
			Success: false,
			Error:   "uid类型转换错误",
		})
		return
	}

	//isasc := true
	// if order == "desc" {
	// 	isasc = false
	// }
	//pagesize := 3

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
			if err := recover(); err != nil {
				fmt.Printf("GetCheckinAll捕获到错误：%v\n", err)
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
			service.Logger.Error("获取打卡参与人数失败", zap.Error(err2))
			cancel()
		}

	}()

	//根据uid查询join表
	wg.Add(1)
	go func() {
		defer wg.Done()

		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("GetUserCheckinJoinByuid捕获到错误：%v\n", err)
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
	var todayWeather WeatherResult
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("GetUserCheckinJoinByuid捕获到错误：%v\n", err)
			}
		}()

		//time.Sleep(1 * time.Second)
		weatherCtx, weatherCancel := context.WithTimeout(ctx, 200*time.Millisecond) //50*time.Millisecond
		defer weatherCancel()

		apiUrl := "http://apis.juhe.cn/simpleWeather/query"
		apiKey := service.Cfg.WeatherAppKey
		data := url.Values{}
		data.Set("key", apiKey)
		data.Set("city", "北京")

		//resp, err := http.Get(apiUrl + "?" + data.Encode())
		req, err := http.NewRequestWithContext(weatherCtx, "GET", apiUrl+"?"+data.Encode(), nil)
		if err != nil {
			service.Logger.Error("NewRequestWithContext err", zap.Error(err))
			weatherErr = err
			return
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			weatherErr = err
			return
		}
		defer resp.Body.Close()

		//var weatherResp map[string]interface{}
		var weather weatherstruct

		//json.NewDecoder(r).Decode(&v)     json.Unmarshal([]byte/*io.ReadAll(resp.Body)*/, &v)两个是一个作用
		//err = json.NewDecoder(resp.Body).Decode(&weatherResp)
		weatherdate, err := io.ReadAll(resp.Body)
		if err != nil {
			service.Logger.Error("ReadAll err", zap.Error(err))
			weatherErr = err
			return
		}

		err = json.Unmarshal(weatherdate, &weather)
		if err != nil {
			service.Logger.Error("Unmarshal err", zap.Error(err))
			weatherErr = err
			return
		}

		todayWeatherInfo := weather.Result.Realtime.Info
		todayTemperature := weather.Result.Realtime.Temperature
		todayWeather = WeatherResult{
			Weather:     todayWeatherInfo,
			Temperature: todayTemperature,
		}

		service.Logger.Info("today weather", zap.String("todayWeather", todayWeatherInfo), zap.String("todaytemperature", todayTemperature))
	}()

	wg.Wait()

	if err1 != nil {
		service.Logger.Error("redis查询失败", zap.String("err:", err1.Error()))
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Error:   "redis查询失败" + err1.Error(),
		})
		return
	}

	if err2 != nil {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Error:   "redis查询失败" + err2.Error(),
		})
		return
	}

	//根据uid获取用户所有的已参与打卡，转换为map[cid]Join，和今日所有的打卡记录转换为map[cid]Record

	if err3 != nil {
		service.Logger.Error("查询参与数据库错误", zap.String("err为", err3.Error()))
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Success: false,
			Error:   "查询参与数据库的错误",
		})
		return
	}
	// if err5 != nil {
	// 	select {
	// 	case <-ctx.Done():
	// 		service.Logger.Error("err", zap.Error(err)) //取消
	// 	case <-weatherCtx.Done():
	// 		c.JSON(http.StatusInternalServerError, model.APIResponse{
	// 			Success: false,
	// 			Error:   "超时",
	// 		})
	// 		service.Logger.Error("err", zap.Error(err))
	// 	default:
	// 		service.Logger.Error("weatherCtx err", zap.Error(err)) //失败
	// 	}
	// 	return
	// }

	if weatherErr != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Success: false,
			Error:   "超时",
		})
		service.Logger.Error("err", zap.Error(err))
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
		service.Logger.Error("查询打卡记录数据库错误", zap.String("err为", err4.Error()))
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Success: false,
			Error:   "查询打卡记录数据库的错误",
		})
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

	datastruct := datastruct{
		Weather: model.WeatherItem{
			Temperature: todayWeather.Temperature,
			Weather:     todayWeather.Weather,
		},
		Slices: responsecheckin,
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Message: "数据库查询排序成功",
		Data:    datastruct,
		//responsecheckin
	})
}
