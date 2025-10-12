package service

import (
	"demo1/internal/app/mydemo/dao/mysql"
	"demo1/internal/app/mydemo/dao/redis"
	"demo1/internal/app/mydemo/model"
	"time"

	"gorm.io/gorm"
)

// 添加checkin到mysql数据
func CreateCheckin(newcheckin *model.Checkin) (result *gorm.DB) {
	return mysql.CreateCheckin(newcheckin)
}

// 保存用户到 Redis 缓存
func SetCheckinToCache(checkin *model.Checkin, ttl time.Duration) error {
	return redis.SetCheckinToCache(checkin, ttl)
}

/*
// 查询checkin的所有redis数据

	func GetCheckinFormCache(id int) (checkin *model.Checkin, err error) {
		return redis.GetCheckinFormCache(id)
	}
*/
//查询全部
func GetCheckinOrderId(page int, pagesize int, isasc bool) (checkins []model.Checkin, err error) {
	return mysql.GetCheckinOrderId(page, pagesize, isasc)
}
