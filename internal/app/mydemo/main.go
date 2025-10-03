package main

import (
	"demo1/internal/app/mydemo/controller"
	"demo1/internal/app/mydemo/service"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
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
		//fmt.Errorf: 包装错误，添加上下文
		err1 := fmt.Errorf("加载配置失败: %w", err)
		//fmt.Sprintf: 构造更详细的错误消息
		strerror := fmt.Sprintf("详细错误%s", err1.Error())
		fmt.Println(strerror)
		log.Fatal(err1)
	}

	// 初始化数据库
	err = service.ServiceInitDB(service.Cfg.Database.Dsn)
	if err != nil {
		err2 := fmt.Errorf("初始化数据库错误%w", err)
		strerror2 := fmt.Sprintf("详细错误%s", err2.Error())
		fmt.Println(strerror2)
		log.Fatal(err2)
	}

	//初始化 Redis
	service.ServiceInitRedis(service.Cfg.Redis.Addr, service.Cfg.Redis.Password, service.Cfg.Redis.DB)

	//创建 Gin 路由引擎
	r := gin.Default()

	// 注册路由
	r.POST("/api/user/add", controller.AddUserHandler)       //绑定路径和函数，当客户端请求路径为""时使用这个函数处理请求
	r.GET("/api/user", controller.GetUserHandler)            //查询id
	r.POST("/api/user/update", controller.UpdateUserHandler) //修改
	r.GET("/api/users/all", controller.GetUsersHandlerAll)   //查询age排序

	r.POST("/api/checkin/add", controller.AddCheckinHandler)   //添加checkin
	r.GET("/api/checkin/all", controller.GetCheckinHandlerAll) //查询Checkin，根据id排序
	// 启动服务器
	log.Println("服务器启动在 :8080 端口")
	log.Fatal(r.Run(":8080"))
}
