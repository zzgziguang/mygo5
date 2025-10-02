package service

import (
	"demo1/internal/app/mydemo/dao/mysql"
	"demo1/internal/app/mydemo/dao/redis"
	"demo1/internal/app/mydemo/model"
	"time"

	"gorm.io/gorm"
)

func CreateCheckin(newcheckin *model.Checkin) (result *gorm.DB) {
	return mysql.CreateCheckin(newcheckin)
}

// 保存用户到 Redis 缓存
func SetCheckinToCache(checkin *model.Checkin, ttl time.Duration) error {
	return redis.SetCheckinToCache(checkin, ttl)
}
