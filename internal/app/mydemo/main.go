package main

import (
	"demo1/internal/app/mydemo/controller"
	"demo1/internal/app/mydemo/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 全局变量：保存 config 配置

func main() {
	var err error

	err = service.LoggerInit()
	if err != nil {
		panic(err)
	}
	//Logger, err = zap.NewDevelopment()
	defer service.SyncLogger()

	err = service.LoadConfig()
	if err != nil {
		service.Logger.Error("LoadConfig err", zap.Error(err))
		panic(err)
	}

	// 初始化数据库
	err = service.ServiceInitDB(service.Cfg.Database.Dsn)
	if err != nil {
		service.Logger.Error("InitDB err", zap.Error(err))
		panic(err)
	}

	//初始化 Redis
	service.ServiceInitRedis(service.Cfg.Redis.Addr, service.Cfg.Redis.Password, service.Cfg.Redis.DB)

	//初始化 kafka
	err = service.ServiceInitKafka()
	if err != nil {
		service.Logger.Error("InitKafka err", zap.Error(err))
		panic(err)
	}
	defer service.Closekafka()

	//创建 Gin 路由引擎
	r := gin.Default()

	// 注册路由
	r.POST("/api/user/add", controller.AddUserHandler)       //绑定路径和函数，当客户端请求路径为""时使用这个函数处理请求
	r.GET("/api/user", controller.GetUserHandler)            //查询id
	r.POST("/api/user/update", controller.UpdateUserHandler) //修改
	r.GET("/api/users/all", controller.GetUsersHandlerAll)   //查询age排序

	r.POST("/api/checkin/add", controller.AddCheckinHandler)   //添加checkin
	r.GET("/api/checkin/all", controller.GetCheckinHandlerAll) //查询Checkin，根据id排序

	//打卡信息
	r.POST("/api/userCheckin/add", controller.AddUserCheckinHandler)
	//获取用户每日打卡
	r.GET("/api/userCheckin/get", controller.GetUserCheckinRecordHandler)

	// 启动服务器
	service.Logger.Info("The server started at port", zap.String("port", "8081"))
	service.Logger.Error("Default error", zap.Error(r.Run(":8081")))
}
