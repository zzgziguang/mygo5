package controller

import (
	"demo1/internal/app/mydemo/model"
	"demo1/internal/app/mydemo/service"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AddUserCheckinHandler(c *gin.Context) {
	struid := c.PostForm("uid")
	if struid == "" {
		MakeApiResponseError(c, CODE_PARAMS_ERROR)
		return
	}

	strcid := c.PostForm("cid")
	if strcid == "" {
		MakeApiResponseError(c, CODE_PARAMS_ERROR)
		return
	}

	uid, err := strconv.Atoi(struid)
	if err != nil {
		service.Logger.Error("uid Atoi", zap.Error(err))
		MakeApiResponseError(c, CODE_PARAMS_ERROR)
		return
	}

	cid, err := strconv.Atoi(strcid)
	if err != nil {
		service.Logger.Error("cid Atoi", zap.Error(err))
		MakeApiResponseError(c, CODE_PARAMS_ERROR)
		return
	}

	createTime := time.Now()
	date := createTime.Year()*10000 + int(createTime.Month())*100 + createTime.Day()

	userCheckinJoinMsg := &model.UserCheckinJoin{}
	userCheckinRecordMsg := &model.UserCheckinRecord{}
	// userCheckin := &model.UserCheckin{
	// 	LastDate:           string,
	// 	DayNum:             int,
	// 	CheckinDayNumMap:   map[int64]int,
	// 	CheckinJoinTimeMap: map[int64]time.Time,
	// 	CheckinLastTimeMap: make(map[int64]time.Time),
	// }

	// endTime, err := service.GetRedisCheckinEndTimeByCid(cid)
	// if err != nil {
	// 	service.Logger.Error("GetRedisCheckinEndTimeByCid err", zap.Error(err))
	// 	MakeApiResponseError(c, CODE_SYS_ERROR)
	// 	return
	// }
	// if endTime < date {
	// 	service.Logger.Error("checkin err", zap.String("err为", "打卡已结束"))
	// 	MakeApiResponseError(c, CODE_SYS_ERROR)
	// 	return
	// }

	// rank, err := service.IncrUserCheckinRecordCountToCache(cid, date)
	// if err != nil {
	// 	service.Logger.Error("IncrUserCheckinRecordCountToCache err", zap.String("err为", err.Error()))
	// 	MakeApiResponseError(c, CODE_SYS_ERROR)
	// 	return
	// }
	// service.Logger.Info("rank uid", zap.Int("rank", int(rank)), zap.Int("uid", uid))
	// MakeApiResponseSuccess(c, map[string]interface{}{
	// 	"rank": rank, //今日打卡名次
	// })
	// return

	//获取用户打卡信息缓存
	userCheckinCache, err := service.HGetUserCheckinFromCache(uid)
	if err != nil {
		service.Logger.Error("HGetUserCheckinFromCache err", zap.String("err为", err.Error()))
		MakeApiResponseError(c, CODE_SYS_ERROR)
		return
	}

	if userCheckinCache == nil {

		userCheckinCache = &model.UserCheckin{
			LastDate: "",
			DayNum:   0,
			CheckinDayNumMap: map[int64]int{
				int64(cid): 0,
			},
			CheckinJoinTimeMap: map[int64]time.Time{
				int64(cid): time.Time{},
			},
			CheckinLastTimeMap: map[int64]time.Time{
				int64(cid): time.Time{},
			},
		}

	}
	//用户是否今日首次打卡
	if userCheckinCache.TodayFristCheckinBool() == false {
		service.Logger.Info("usertodayfirstcheckin", zap.String("usertodayfirstcheckin", "用户今日首次打卡"))

		//上次打卡日期
		err = service.HSetUserCheckinLastDateToCache(uid, strconv.Itoa(date))
		if err != nil {
			service.Logger.Error("HSetUserCheckinLastDateToCache err", zap.String("err为", err.Error()))
			MakeApiResponseError(c, CODE_SYS_ERROR)
			return
		}

	}
	//cid是否首次打卡
	if userCheckinCache.FristJoinCheckinBool(int64(cid)) == false {
		service.Logger.Info("userfirstcheckincid", zap.String("userfirstcheckincid", "用户首次打卡cid"))

		//添加jointime缓存
		err = service.HSetUserCheckinJoinTimeToCache(uid, cid, createTime.Unix())
		if err != nil {
			service.Logger.Error("AddUserCheckinJoin err", zap.Error(err))
			MakeApiResponseErrorDefault(c)
			return
		}

		userCheckinJoinMsg = &model.UserCheckinJoin{
			Uid:      uid,
			Cid:      cid,
			JoinTime: &createTime,
			CreateAt: &createTime,
			UpdateAt: &createTime,
			Status:   model.JoinStatusNormal,
		}

		partition, offset, err := service.ProduceKafkaUserCheckinJoinMessage(userCheckinJoinMsg)
		if err != nil {
			service.Logger.Error("ProduceKafkaUserCheckinJoinMessage err", zap.Error(err))
			MakeApiResponseErrorDefault(c)
			return
		}
		service.Logger.Debug("ProduceKafkaUserCheckinJoinMessage", zap.Any("partition", partition), zap.Any("offset", offset))

	}

	//cid是否今日首次打卡
	if userCheckinCache.FristDateJoinCheckinBool(int64(cid)) == true {
		service.Logger.Info("usertodayfirstcheckincid", zap.String("userfirstcheckincid", "用户今日首次打卡cid"))
		//上次打卡时间
		err = service.HSetUserCheckinLastTimeToCache(uid, cid, createTime.Unix())
		if err != nil {
			service.Logger.Error("HSetUserCheckinLastTimeToCache err", zap.String("err为", err.Error()))
			MakeApiResponseError(c, CODE_SYS_ERROR)
			return
		}

		//某打卡天数
		err = service.HSetUserCheckinDateNumToCache(uid, cid)
		if err != nil {
			service.Logger.Error("HSetUserCheckinDateNumToCache err", zap.String("err为", err.Error()))
			MakeApiResponseError(c, CODE_SYS_ERROR)
			return
		}

		//用户打卡总天数
		err = service.HSetUserCheckinDayNumToCache(uid)
		if err != nil {
			service.Logger.Error("HSetUserCheckinDayNumToCache err", zap.String("err为", err.Error()))
			MakeApiResponseError(c, CODE_SYS_ERROR)
			return
		}
	}

	userCheckinRecordMsg = &model.UserCheckinRecord{
		Uid:      uid,
		Cid:      cid,
		Date:     date,
		CreateAt: &createTime,
		UpdateAt: &createTime,
		Status:   model.JoinStatusNormal,
	}

	partition, offset, err := service.ProduceKafkaUserCheckinRecordMessage(userCheckinRecordMsg)
	if err != nil {
		service.Logger.Error("ProduceKafkaUserCheckinRecordMessage err", zap.Error(err))
		MakeApiResponseErrorDefault(c)
		return
	}
	service.Logger.Debug("ProduceKafkaUserCheckinRecordMessage", zap.Any("partition", partition), zap.Any("offset", offset))

	rank, err := service.IncrUserCheckinRecordCountToCache(cid, date)
	if err != nil {
		service.Logger.Error("IncrUserCheckinRecordCountToCache err", zap.String("err为", err.Error()))
		MakeApiResponseError(c, CODE_SYS_ERROR)
		return
	}

	// err = service.AddUserCheckinRecord(userCheckinRecordMsg)
	// if err != nil {
	// 	service.Logger.Error("AddUserCheckinRecord err", zap.Error(err))
	// 	MakeApiResponseError(c, CODE_SYS_ERROR)
	// 	return
	// }

	MakeApiResponseSuccess(c, map[string]interface{}{
		"userCheckin":       userCheckinCache,     //用户打卡数据
		"userCheckinJoin":   userCheckinJoinMsg,   //参与表中数据
		"userCheckinRecord": userCheckinRecordMsg, //打卡记录表数据
		"rank":              rank,                 //今日打卡名次
	})

	// return

	// //get join from cache
	// userCheckinJoin, err := service.GetUserCheckinJoinFromCache(uid, cid)
	// if err != nil {
	// 	service.Logger.Error("GetUserCheckinJoinFromCache", zap.Error(err))
	// 	MakeApiResponseErrorDefault(c)
	// 	return
	// } else {
	// 	//无错误
	// 	//返回return
	// }
	// if userCheckinJoin != nil {
	// 	//缓存中有join数据
	// 	//直接使用缓存的join数据
	// 	service.Logger.Debug("userCheckinJoin!=nil", zap.String("userCheckinJoin=", fmt.Sprintf("%v", userCheckinJoin)))
	// } else {
	// 	//缓存中没有join数据
	// 	//从db中取join数据
	// 	userCheckinJoin, err = service.GetUserCheckinJoin(uid, cid)
	// 	if err != nil {
	// 		service.Logger.Error("userCheckinJoin!=nil", zap.Error(err))
	// 		MakeApiResponseErrorDefault(c)
	// 		return
	// 	}

	// 	if userCheckinJoin == nil {
	// 		//未参与
	// 		//参与，写db
	// 		joinTime := time.Now()
	// 		userCheckinJoin = &model.UserCheckinJoin{
	// 			Uid:      uid,
	// 			Cid:      cid,
	// 			JoinTime: &joinTime,
	// 			CreateAt: &joinTime,
	// 			UpdateAt: &joinTime,
	// 			Status:   model.JoinStatusNormal,
	// 		}

	// 		err = service.AddUserCheckinJoin(userCheckinJoin)
	// 		if err != nil {
	// 			service.Logger.Error("AddUserCheckinJoin err", zap.Error(err))
	// 			MakeApiResponseErrorDefault(c)
	// 			return
	// 		}

	// 		service.Logger.Debug("AddUserCheckinJoin", zap.String("userCheckinJoin", fmt.Sprintf("%V", userCheckinJoin)))

	// 		//添加kafka生产者
	// 		msg := model.CheckInMsg{
	// 			Uid:       uid,
	// 			Cid:       cid,
	// 			Timestamp: time.Now().Unix(),
	// 			Msg:       fmt.Sprintf("恭喜%d打卡%d成功", uid, cid),
	// 		}

	// 		// 序列化为 json
	// 		value, err := json.Marshal(msg)
	// 		if err != nil {
	// 			service.Logger.Error("Marshal Error", zap.Error(err))
	// 			MakeApiResponseErrorDefault(c)
	// 			return
	// 		}

	// 		for i := 1; i <= 10; i++ {
	// 			partition, offset, err := service.ProducerSend(value)
	// 			if err != nil {
	// 				service.Logger.Error("ProducerSend err", zap.Error(err))
	// 				MakeApiResponseErrorDefault(c)
	// 				return
	// 			}
	// 			service.Logger.Debug("ProducerSend", zap.Any("partition", partition), zap.Any("offset", offset))
	// 		}

	// 		//获取checkin表的id=cid，看看有没有这个打卡
	// 		checkin, err := service.GetCheckinBycid(cid)
	// 		if err != nil {
	// 			service.Logger.Error("GetCheckinBycid err", zap.Error(err))
	// 			MakeApiResponseErrorDefault(c)
	// 			return
	// 		}

	// 		//有打卡，就更新参与打卡人数
	// 		//更新
	// 		err = service.UpdateCheckinJoinNum(cid, checkin)
	// 		if err != nil {
	// 			service.Logger.Error("UpdateCheckinJoinNum err", zap.Error(err))
	// 			MakeApiResponseErrorDefault(c)
	// 			return
	// 		}

	// 		//将更新后打卡人数添加到zset
	// 		err = service.ZaddCheckinJoinNum(cid, checkin)
	// 		if err != nil {
	// 			service.Logger.Error("ZaddCheckinJoinNum err", zap.Error(err))
	// 			MakeApiResponseErrorDefault(c)
	// 			return
	// 		}

	// 	} else {
	// 		service.Logger.Debug("GetUserCheckinJoin", zap.String("userCheckinJoin", fmt.Sprintf("%V", userCheckinJoin)))
	// 	}
	// }

	// //get record
	// userCheckinRecord, err := service.GetUserCheckinRecord(uid, cid, date)
	// if err != nil {
	// 	service.Logger.Error("GetUserCheckinRecord err", zap.Error(err))
	// 	MakeApiResponseErrorDefault(c)
	// 	return
	// }

	// if userCheckinRecord != nil {
	// 	//今日已打卡,返回已打卡
	// 	MakeApiResponseSuccess(c, map[string]interface{}{
	// 		"userCheckinJoin":   userCheckinJoin,   //参与表数据
	// 		"userCheckinRecord": userCheckinRecord, //打卡记录表数据
	// 	})
	// 	return
	// } else {

	// 	userCheckinRecord = &model.UserCheckinRecord{
	// 		Uid:      uid,
	// 		Cid:      cid,
	// 		Date:     date,
	// 		CreateAt: &createTime,
	// 		UpdateAt: &createTime,
	// 		Status:   model.JoinStatusNormal,
	// 	}

	// 	err = service.AddUserCheckinRecord(userCheckinRecord)
	// 	if err != nil {
	// 		service.Logger.Error("AddUserCheckinRecord err", zap.Error(err))
	// 		MakeApiResponseError(c, CODE_SYS_ERROR)
	// 		return
	// 	}

	// 	// rank, err := service.GetUserCheckinRecordByCount(cid, date, recordTime)
	// 	// if err != nil {
	// 	// 	service.Logger.Error("GetUserCheckinRecordByCount err", zap.String("err为", err.Error()))
	// 	// 	MakeApiResponseError(c, CODE_SYS_ERROR)
	// 	// 	return
	// 	// }

	// 	// err = service.ZaddUserCheckinRecordCountToCache(cid, uid, recordTime, date)
	// 	// if err != nil {
	// 	// 	service.Logger.Error("ZaddUserCheckinRecordCountToCache err", zap.String("err为", err.Error()))
	// 	// 	MakeApiResponseError(c, CODE_SYS_ERROR)
	// 	// 	return
	// 	// }

	// 	// rank, err := service.ZrankUserCheckinRecordCountToCache(cid, uid, date)
	// 	// if err != nil {
	// 	// 	service.Logger.Error("ZrankUserCheckinRecordCountToCache err", zap.String("err为", err.Error()))
	// 	// 	MakeApiResponseError(c, CODE_SYS_ERROR)
	// 	// 	return
	// 	// }

	// 	rank, err := service.IncrUserCheckinRecordCountToCache(cid, date)
	// 	if err != nil {
	// 		service.Logger.Error("IncrUserCheckinRecordCountToCache err", zap.String("err为", err.Error()))
	// 		MakeApiResponseError(c, CODE_SYS_ERROR)
	// 		return
	// 	}

	// 	MakeApiResponseSuccess(c, map[string]interface{}{
	// 		"userCheckin":       userCheckinCache,  //用户打卡数据
	// 		"userCheckinJoin":   userCheckinJoin,   //参与表中数据
	// 		"userCheckinRecord": userCheckinRecord, //打卡记录表数据
	// 		"rank":              rank,              //今日打卡名次
	// 	})
	// }
}

// list
func GetUserCheckinRecordHandler(c *gin.Context) {
	struid := c.PostForm("uid")
	if struid == "" {
		MakeApiResponseError(c, CODE_PARAMS_ERROR)
		return
	}

	strcid := c.PostForm("cid")
	if strcid == "" {
		MakeApiResponseError(c, CODE_PARAMS_ERROR)
		return
	}

	uid, err := strconv.Atoi(struid)
	if err != nil {
		service.Logger.Error("Atoiuid err", zap.Error(err))
		MakeApiResponseError(c, CODE_PARAMS_ERROR)
		return
	}

	cid, err := strconv.Atoi(strcid)
	if err != nil {
		service.Logger.Error("Atoicid err", zap.Error(err))
		MakeApiResponseError(c, CODE_PARAMS_ERROR)
		return
	}

	recordTime := time.Now()
	date := recordTime.Year()*10000 + int(recordTime.Month())*100 + recordTime.Day()

	//设置升序
	isasc := true
	//获取全部打卡列表
	//select * from record表 where uid=1 cid=1
	userCheckinRecordList, err := service.GetUserCheckinRecordList(uid, cid, isasc)
	if err != nil {
		service.Logger.Error("GetUserCheckinRecordList err", zap.Error(err))
		MakeApiResponseErrorDefault(c)
		return
	}

	var strUserCheckinRecord string
	var userCheckinRecord interface{}
	if len(userCheckinRecordList) > 0 { //优先使用len判断slice是否有数据，比只判断=nil更健壮
		//有record数据
		for _, v := range userCheckinRecordList {
			if v.Date == date {
				//今天已打卡
				strUserCheckinRecord = "用户今日已打卡"
				userCheckinRecord = v
				service.Logger.Debug("strUserCheckinRecord", zap.String("v", fmt.Sprintf("%V", v)))
			} else {
				//今天没有打卡
				strUserCheckinRecord = "用户今日未打卡"
				userCheckinRecord = map[string]string{}
				service.Logger.Debug("strUserCheckinRecord", zap.String("v", fmt.Sprintf("%V", v)))
			}
		}
	}

	//未参与
	//查询join表
	userCheckinJoin, err := service.GetUserCheckinJoin(uid, cid)
	if err != nil {
		service.Logger.Error("GetUserCheckinJoin err", zap.Error(err))
		MakeApiResponse(c, 1, "查询参与的错误"+err.Error())
		MakeApiResponseErrorDefault(c)
		return
	}

	//一个数据只可以使用=nil判断是否有数据
	if userCheckinJoin == nil {
		service.Logger.Error("userCheckinJoin err", zap.String("userCheckinJoin", "userCheckinJoin为空")) // 只要json返回的是错误，就记录日志
		MakeApiResponseError(c, CODE_SYS_ERROR)
		return
	} else {
		MakeApiResponseSuccess(c, map[string]interface{}{
			"strUserCheckinRecord": strUserCheckinRecord, //今日有没有打卡
			"userCheckinRecord":    userCheckinRecord,    //打卡数据
			"userCheckinJoin":      userCheckinJoin,      //参与表中数据
		})
		return
	}
}
