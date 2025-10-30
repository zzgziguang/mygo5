package controller

import (
	"demo1/internal/app/mydemo/model"
	"demo1/internal/app/mydemo/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AddUserCheckinHandler(c *gin.Context) {
	// AddUserCheckinHandler2(c)
	// return

	struid := c.PostForm("uid")
	if struid == "" {
		c.JSON(http.StatusNotFound, model.APIResponse{
			Success: false,
			Error:   "uid不能为空",
		})
		return
	}

	strcid := c.PostForm("cid")
	if strcid == "" {
		c.JSON(http.StatusNotFound, model.APIResponse{
			Success: false,
			Error:   "cid不能为空",
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

	cid, err := strconv.Atoi(strcid)
	if err != nil {
		c.JSON(http.StatusNotFound, model.APIResponse{
			Success: false,
			Error:   "cid类型转换错误",
		})
		return
	}

	msg := model.CheckInMsg{
		Uid:       uid,
		Cid:       cid,
		Timestamp: time.Now().Unix(),
		Msg:       fmt.Sprintf("恭喜%d打卡%d成功", uid, cid),
	}

	// 序列化为 JSON
	value, err := json.Marshal(msg)
	if err != nil {
		service.Logger.Error("Marshal Error", zap.Error(err))
		c.JSON(http.StatusNotFound, model.APIResponse{
			Success: false,
			Error:   "json 转换错误err为" + err.Error(),
		})
		return
	}

	for i := 1; i <= 10; i++ {
		partition, offset, err := service.ProducerSend(value)
		if err != nil {
			service.Logger.Error("ProducerSend", zap.String("err", err.Error()))
			c.JSON(http.StatusNotFound, model.APIResponse{
				Success: false,
				Error:   "错误err为" + err.Error(),
			})
			return
		}
		service.Logger.Debug("ProducerSend", zap.Any("partition", partition), zap.Any("offset", offset))
	}

	recordTime := time.Now()
	date := recordTime.Year()*10000 + int(recordTime.Month())*100 + recordTime.Day()

	//get join from cache
	userCheckinJoin, err := service.GetUserCheckinJoinFromCache(uid, cid)
	if err != nil {
		// 有错误
		//if err redisn il
		service.Logger.Error("getUserCheckinJoinFromCacheErr", zap.String("err=", err.Error()))
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Success: false,
			Error:   "缓存获取参与错误",
		})
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
			service.Logger.Error("userCheckinJoin!=nil", zap.String("userCheckinJoin=", fmt.Sprintf("%v", userCheckinJoin)))
			c.JSON(http.StatusInternalServerError, model.APIResponse{
				Success: false,
				Error:   "查询参与数据库的错误",
			})
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
				service.Logger.Error("json表添加错误", zap.String("err为", err.Error()))
				c.JSON(http.StatusNotFound, model.APIResponse{
					Success: false,
					Error:   "错误err为" + err.Error(),
				})
				return
			}
			service.Logger.Debug("添加参与数据库成功", zap.String("userCheckinJoin", fmt.Sprintf("%V", userCheckinJoin)))
			// 添加kafka生产者
			// msg := model.CheckInMsg{
			// 	Uid:       uid,
			// 	Cid:       cid,
			// 	Timestamp: time.Now().Unix(),
			// 	Msg:       fmt.Sprintf("恭喜%d打卡%d成功", uid, cid),
			// }

			// // 序列化为 json
			// value, err := json.Marshal(msg)
			// if err != nil {
			// 	service.Logger.Error("Marshal Error", zap.Error(err))
			// 	c.JSON(http.StatusNotFound, model.APIResponse{
			// 		Success: false,
			// 		Error:   "json 转换错误err为" + err.Error(),
			// 	})
			// 	return
			// }
			// strValue := string(value)

			// for i := 1; i <= 10; i++ {
			// 	partition, offset, err := service.ProducerSend(strValue)
			// 	if err != nil {
			// 		service.Logger.Error("ProducerSend", zap.String("err", err.Error()))
			// 		c.JSON(http.StatusNotFound, model.APIResponse{
			// 			Success: false,
			// 			Error:   "错误err为" + err.Error(),
			// 		})
			// 		return
			// 	}
			// 	service.Logger.Debug("ProducerSend", zap.Any("partition", partition), zap.Any("offset", offset))
			// }
			//获取checkin表的id=cid，看看有没有这个打卡
			checkin, err := service.GetCheckinBycid(cid)
			if err != nil {
				c.JSON(http.StatusInternalServerError, model.APIResponse{
					Success: false,
					Error:   "查询checkin表错误",
				})
				return
			}
			//有打卡，就更新参与打卡人数
			//更新
			err = service.UpdateCheckinJoinNum(cid, checkin)
			if err != nil {
				c.JSON(http.StatusInternalServerError, model.APIResponse{
					Success: false,
					Error:   "更新参与人数失败",
				})
				return
			}

			//将更新后打卡人数添加到zset
			err = service.ZaddCheckinJoinNum(cid, checkin)
			if err != nil {
				c.JSON(http.StatusInternalServerError, model.APIResponse{
					Success: false,
					Error:   "添加参与人数缓存失败",
				})
				return
			}

			// //更新参与打卡的权重
			// err = service.UpdateCheckinWeight(cid, checkin)
			// if err != nil {
			// 	c.JSON(http.StatusInternalServerError, model.APIResponse{
			// 		Success: false,
			// 		Error:   "更新参与人数失败",
			// 	})
			// 	return
			// }

		} else {
			// set cache join，首次join不用缓存，非首次join才用添加
			// err = service.SetUserCheckinJoinToCache(userCheckinJoin)
			// if err != nil {
			// 	c.JSON(http.StatusInternalServerError, model.APIResponse{
			// 		Success: false,
			// 		Error:   "添加用户参与缓存失败",
			// 	})
			// 	return
			// }
			service.Logger.Debug("查询参与数据库成功", zap.String("userCheckinJoin", fmt.Sprintf("%V", userCheckinJoin)))
		}

	}

	//get record
	userCheckinRecord, err := service.GetUserCheckinRecord(uid, cid, date)
	if err != nil {
		service.Logger.Error("查询打卡错误", zap.String("err为", err.Error()))
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Success: false,
			Error:   "查询打卡错误",
		})
		return
	}
	if userCheckinRecord != nil {
		//今日已打卡
		//返回已打卡
		c.JSON(http.StatusOK, model.APIResponse{
			Success: true,
			Message: "用户今日已打卡",
			Data: map[string]interface{}{
				"userCheckinJoin":   userCheckinJoin,   //参与表数据
				"userCheckinRecord": userCheckinRecord, //打卡记录表数据
			},
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
			service.Logger.Error("打卡错误", zap.String("err为", err.Error()))
			c.JSON(http.StatusInternalServerError, model.APIResponse{
				Success: false,
				Error:   "打卡错误",
			})
			return
		}
		a, err := json.Marshal(userCheckinJoin)
		if err != nil {
			return
		}
		fmt.Println(string(a))
		c.JSON(http.StatusOK, model.APIResponse{
			Success: true,
			Message: "用户打卡成功",
			Data: map[string]interface{}{
				"userCheckinJoin":   userCheckinJoin,   //参与表中数据
				"userCheckinRecord": userCheckinRecord, //打卡记录表数据
			},
		})

	}

}

// list
func GetUserCheckinRecordHandler(c *gin.Context) {

	struid := c.PostForm("uid")
	if struid == "" {
		c.JSON(http.StatusNotFound, model.APIResponse{
			Success: false,
			Error:   "uid不能为空",
		})
		return
	}

	strcid := c.PostForm("cid")
	if strcid == "" {
		c.JSON(http.StatusNotFound, model.APIResponse{
			Success: false,
			Error:   "cid不能为空",
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
	cid, err := strconv.Atoi(strcid)
	if err != nil {
		c.JSON(http.StatusNotFound, model.APIResponse{
			Success: false,
			Error:   "cid类型转换错误",
		})
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
		service.Logger.Error("查询全部打卡记录错误", zap.String("err为", err.Error()))
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Success: false,
			Error:   "查询全部打卡记录错误",
		})
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
				service.Logger.Debug("添加参与数据库成功", zap.String("v", fmt.Sprintf("%V", v)))
			} else {
				//今天没有打卡
				strUserCheckinRecord = "用户今日未打卡"
				userCheckinRecord = map[string]string{}
				service.Logger.Debug("添加参与数据库成功", zap.String("v", fmt.Sprintf("%V", v)))
			}
		}

	}
	//未参与
	//查询join表
	userCheckinJoin, err := service.GetUserCheckinJoin(uid, cid)
	if err != nil {
		service.Logger.Error("查询参与错误", zap.String("err为", err.Error()))
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Success: false,
			Error:   "查询参与的错误",
		})
		return
	}

	if userCheckinJoin == nil { //一个数据只可以使用=nil判断是否有数据
		c.JSON(http.StatusOK, model.APIResponse{
			Success: true,
			Message: "用户未参与打卡",
			Data:    map[string]string{},
		})
		return
	} else {
		c.JSON(http.StatusOK, model.APIResponse{
			Success: true,
			Message: "用户已参与打卡",
			Data: map[string]interface{}{
				"strUserCheckinRecord": strUserCheckinRecord, //今日有没有打卡
				"userCheckinRecord":    userCheckinRecord,    //打卡数据
				"userCheckinJoin":      userCheckinJoin,      //参与表中数据

			},
		})
	}

}
