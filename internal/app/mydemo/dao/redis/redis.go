package redis

import (
	"log"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func DaoInitRedis(addr, password string, db int) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,     // Redis 地址
		Password: password, // 无密码
		DB:       db,       // 默认 DB
	})

	// 测试连接
	/*
		_, err := RedisClient.Ping(Ctx).Result()
		if err != nil {
			log.Fatal("Redis 连接失败: ", err)
		}
	*/
	log.Println("Redis 连接成功")
	return
}
