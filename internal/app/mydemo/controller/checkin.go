package controller

import (
	"demo1/internal/app/mydemo/model"
	"demo1/internal/app/mydemo/service"
	"net/http"
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
