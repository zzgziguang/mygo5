package controller

import (
	"demo1/internal/app/mydemo/model"
	"demo1/internal/app/mydemo/service"
	"encoding/json"
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

	recordTime := time.Now()
	date := recordTime.Year()*10000 + int(recordTime.Month())*100 + recordTime.Day()

	rank, err := service.IncrUserCheckinRecordCountToCache(cid, date)
	if err != nil {
		service.Logger.Error("IncrUserCheckinRecordCountToCache err", zap.String("err为", err.Error()))
		MakeApiResponseError(c, CODE_SYS_ERROR)
		return
	}

	MakeApiResponseSuccess(c, map[string]interface{}{
		"rank": rank, //今日打卡名次
	})
	return

	//get join from cache
	userCheckinJoin, err := service.GetUserCheckinJoinFromCache(uid, cid)
	if err != nil {
		service.Logger.Error("GetUserCheckinJoinFromCache", zap.Error(err))
		MakeApiResponseErrorDefault(c)
		return
	} else {
		//无错误
		//返回return
	}
	if userCheckinJoin != nil {
		//缓存中有join数据
		//直接使用缓存的join数据
		service.Logger.Debug("userCheckinJoin!=nil", zap.String("userCheckinJoin=", fmt.Sprintf("%v", userCheckinJoin)))
	} else {
		//缓存中没有join数据
		//从db中取join数据
		userCheckinJoin, err = service.GetUserCheckinJoin(uid, cid)
		if err != nil {
			service.Logger.Error("userCheckinJoin!=nil", zap.Error(err))
			MakeApiResponseErrorDefault(c)
			return
		}

		if userCheckinJoin == nil {
			//未参与
			//参与，写db
			joinTime := time.Now()
			userCheckinJoin = &model.UserCheckinJoin{
				Uid:      uid,
				Cid:      cid,
				JoinTime: &joinTime,
				CreateAt: &joinTime,
				UpdateAt: &joinTime,
				Status:   model.JoinStatusNormal,
			}

			err = service.AddUserCheckinJoin(userCheckinJoin)
			if err != nil {
				service.Logger.Error("AddUserCheckinJoin err", zap.Error(err))
				MakeApiResponseErrorDefault(c)
				return
			}

			service.Logger.Debug("AddUserCheckinJoin", zap.String("userCheckinJoin", fmt.Sprintf("%V", userCheckinJoin)))

			//添加kafka生产者
			msg := model.CheckInMsg{
				Uid:       uid,
				Cid:       cid,
				Timestamp: time.Now().Unix(),
				Msg:       fmt.Sprintf("恭喜%d打卡%d成功", uid, cid),
			}

			// 序列化为 json
			value, err := json.Marshal(msg)
			if err != nil {
				service.Logger.Error("Marshal Error", zap.Error(err))
				MakeApiResponseErrorDefault(c)
				return
			}

			for i := 1; i <= 10; i++ {
				partition, offset, err := service.ProducerSend(value)
				if err != nil {
					service.Logger.Error("ProducerSend err", zap.Error(err))
					MakeApiResponseErrorDefault(c)
					return
				}
				service.Logger.Debug("ProducerSend", zap.Any("partition", partition), zap.Any("offset", offset))
			}

			//获取checkin表的id=cid，看看有没有这个打卡
			checkin, err := service.GetCheckinBycid(cid)
			if err != nil {
				service.Logger.Error("GetCheckinBycid err", zap.Error(err))
				MakeApiResponseErrorDefault(c)
				return
			}

			//有打卡，就更新参与打卡人数
			//更新
			err = service.UpdateCheckinJoinNum(cid, checkin)
			if err != nil {
				service.Logger.Error("UpdateCheckinJoinNum err", zap.Error(err))
				MakeApiResponseErrorDefault(c)
				return
			}

			//将更新后打卡人数添加到zset
			err = service.ZaddCheckinJoinNum(cid, checkin)
			if err != nil {
				service.Logger.Error("ZaddCheckinJoinNum err", zap.Error(err))
				MakeApiResponseErrorDefault(c)
				return
			}

		} else {
			service.Logger.Debug("GetUserCheckinJoin", zap.String("userCheckinJoin", fmt.Sprintf("%V", userCheckinJoin)))
		}
	}

	//get record
	userCheckinRecord, err := service.GetUserCheckinRecord(uid, cid, date)
	if err != nil {
		service.Logger.Error("GetUserCheckinRecord err", zap.Error(err))
		MakeApiResponseErrorDefault(c)
		return
	}

	if userCheckinRecord != nil {
		//今日已打卡,返回已打卡
		MakeApiResponseSuccess(c, map[string]interface{}{
			"userCheckinJoin":   userCheckinJoin,   //参与表数据
			"userCheckinRecord": userCheckinRecord, //打卡记录表数据
		})
		return
	} else {
		userCheckinRecord = &model.UserCheckinRecord{
			Uid:      uid,
			Cid:      cid,
			Date:     date,
			CreateAt: &recordTime,
			UpdateAt: &recordTime,
			Status:   model.JoinStatusNormal,
		}

		err = service.AddUserCheckinRecord(userCheckinRecord)
		if err != nil {
			service.Logger.Error("AddUserCheckinRecord err", zap.String("err为", err.Error()))
			MakeApiResponseError(c, CODE_SYS_ERROR)
			return
		}

		// rank, err := service.GetUserCheckinRecordByCount(cid, date, recordTime)
		// if err != nil {
		// 	service.Logger.Error("GetUserCheckinRecordByCount err", zap.String("err为", err.Error()))
		// 	MakeApiResponseError(c, CODE_SYS_ERROR)
		// 	return
		// }

		// err = service.ZaddUserCheckinRecordCountToCache(cid, uid, recordTime, date)
		// if err != nil {
		// 	service.Logger.Error("ZaddUserCheckinRecordCountToCache err", zap.String("err为", err.Error()))
		// 	MakeApiResponseError(c, CODE_SYS_ERROR)
		// 	return
		// }

		// rank, err := service.ZrankUserCheckinRecordCountToCache(cid, uid, date)
		// if err != nil {
		// 	service.Logger.Error("ZrankUserCheckinRecordCountToCache err", zap.String("err为", err.Error()))
		// 	MakeApiResponseError(c, CODE_SYS_ERROR)
		// 	return
		// }

		rank, err := service.IncrUserCheckinRecordCountToCache(cid, date)
		if err != nil {
			service.Logger.Error("IncrUserCheckinRecordCountToCache err", zap.String("err为", err.Error()))
			MakeApiResponseError(c, CODE_SYS_ERROR)
			return
		}

		MakeApiResponseSuccess(c, map[string]interface{}{
			"userCheckinJoin":   userCheckinJoin,   //参与表中数据
			"userCheckinRecord": userCheckinRecord, //打卡记录表数据
			"rank":              rank,              //今日打卡名次
		})
	}
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
		service.Logger.Error("") //todo 只要json返回的是错误，就记录日志
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
