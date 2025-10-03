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
	// 添加成功后写入 Redis 缓存
	var ttl = time.Duration(service.Cfg.Redis.CacheTTL) * time.Second
	err := service.SetCheckinToCache(newCheckin, ttl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.APIResponse{
			Success: false,
			Error:   "写入redis失败: " + err.Error(),
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

	isasc := true
	if order == "desc" {
		isasc = false
	}
	pagesize := 3
	result, err := service.GetCheckinOrderId(page, pagesize, isasc)
	if err != nil {
		service.Logger.Error("数据库查询失败", zap.String("err:", err.Error()))
		c.JSON(http.StatusBadRequest, model.APIResponse{
			Success: false,
			Error:   "数据库查询失败" + err.Error(),
		})
		return
	}

	responsecheckin := make([]model.ResponseCheckinItem, 0)
	//var responsecheckin []model.ResponseCheckinItem=make([]model.ResponseCheckinItem, 3,3)
	//responsecheckin := make([]model.ResponseCheckinItem,0)
	for _, v := range result {
		checkinre := model.ResponseCheckinItem{
			Id:            v.Id,
			Title:         v.Title,
			CreateAt:      v.CreateAt.Format("2006年01月02日 15点04分05秒"),
			UpdateAt:      v.UpdateAt.Format("2006年01月02日 15点04分05秒"),
			CheckinStatus: v.CheckinStatus,
		}
		responsecheckin = append(responsecheckin, checkinre)
	}

	c.JSON(http.StatusOK, model.APIResponse{
		Success: true,
		Message: "数据库查询排序成功",
		Data:    responsecheckin,
	})
}
