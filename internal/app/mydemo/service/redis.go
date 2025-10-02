package service

import "demo1/internal/app/mydemo/dao/redis"

// 初始化 Redis 客户端
func ServiceInitRedis(addr, password string, db int) {
	redis.DaoInitRedis(addr, password, db)
}
