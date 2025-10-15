package controller

import (
	"demo1/internal/app/mydemo/model"
	"demo1/internal/app/mydemo/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AddCheckinHandler(c *gin.Context) {
	title := c.PostForm("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Error:   "标题不能为空",
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

	isasc := true
	if order == "desc" {
		isasc = false
	}
	pagesize := 3

	//获取全部
	//checkinData，checkinList,checkinSlice,checkinResult,checkinRes
	checkinSlice, err := service.GetCheckinOrderJoinnumber(page, pagesize, isasc)
	if err != nil {
		service.Logger.Error("redis查询失败", zap.String("err:", err.Error()))
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Error:   "redis查询失败" + err.Error(),
		})
		return
	}

	time := time.Now()
	date := time.Year()*10000 + int(time.Month())*100 + time.Day()

	//获取zset缓存
	// zrank, err = service.GetZsetCheckinNum(cid)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, model.APIResponse{
	// 		Success: false,
	// 		Error:   "redis查询失败" + err.Error(),
	// 	})
	// 	return
	// }

	//根据uid获取用户所有的已参与打卡，转换为map[cid]Join，和今日所有的打卡记录转换为map[cid]Record
	//根据uid查询join表
	joinSlice, err := service.GetUserCheckinJoinByuid(uid)
	if err != nil {
		service.Logger.Error("查询参与数据库错误", zap.String("err为", err.Error()))
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Success: false,
			Error:   "查询参与数据库的错误",
		})
		return
	}
	userCheckinJoinMap := make(map[int]*model.UserCheckinJoin, 0)
	for _, userCheckinJoin := range joinSlice {
		userCheckinJoinMap[userCheckinJoin.Cid] = userCheckinJoin
	}

	recordSlice, err := service.GetUserCheckinRecordByUidDate(uid, date)
	if err != nil {
		service.Logger.Error("查询参与数据库错误", zap.String("err为", err.Error()))
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Success: false,
			Error:   "查询参与数据库的错误",
		})
		return
	}
	userCheckinRecordMap := make(map[int]*model.UserCheckinRecord, 0)
	for _, userCheckinRecord := range recordSlice {
		userCheckinRecordMap[userCheckinRecord.Cid] = &userCheckinRecord
	}

	responsecheckin := make([]model.ResponseCheckinItem, 0)
	//var responsecheckin []model.ResponseCheckinItem=make([]model.ResponseCheckinItem, 3,3)
	//responsecheckin := make([]model.ResponseCheckinItem,0)

	for _, v := range checkinSlice {
		cid := v.Id
		joinBool := false
		recordBool := false
		_, ok := userCheckinJoinMap[cid]
		if ok {
			joinBool = true
		}
		_, ok = userCheckinRecordMap[cid]
		if ok {
			recordBool = true
		}

		checkinre := model.ResponseCheckinItem{
			Id:            v.Id,
			Title:         v.Title,
			CreateAt:      v.CreateAt.Format("2006年01月02日 15点04分05秒"),
			UpdateAt:      v.UpdateAt.Format("2006年01月02日 15点04分05秒"),
			CheckinStatus: v.CheckinStatus,
			JoinBool:      joinBool,         //是否参与
			RecordBool:    recordBool,       //是否打卡
			JoinNumber:    int64(v.JoinNum), //参与人数
			//是否已参与 _, ok := map[cid]
			//今日是否已打卡
		}
		responsecheckin = append(responsecheckin, checkinre)
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Message: "数据库查询排序成功",
		Data:    responsecheckin,
	})
}
