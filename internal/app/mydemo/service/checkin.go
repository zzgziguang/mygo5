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

// 查询全部id
func GetCheckinOrderId(page int, pagesize int, isasc bool) (checkins []model.Checkin, err error) {
	return mysql.GetCheckinOrderId(page, pagesize, isasc)
}

// 查询全部
func GetCheckinAll() (checkins []model.Checkin, err error) {
	return mysql.GetCheckinAll()
}

// 查询全部joinnumber
func GetCheckinOrderJoinnumber(page int, pagesize int, isasc bool) (checkins []model.Checkin, err error) {
	return mysql.GetCheckinOrderJoinnumber(page, pagesize, isasc)
}

// 查询id=cid
func GetCheckinBycid(cid int) (checkin *model.Checkin, err error) {

	return mysql.GetCheckinBycid(cid)
}

// 更新joinnumber
func UpdateCheckinJoinNum(cid int, checkin *model.Checkin) (err error) {
	return mysql.UpdateCheckinJoinNum(cid, checkin)
}

// 更新weight
func UpdateCheckinWeight(newCheckin *model.Checkin, weight int) (err error) {
	return mysql.UpdateCheckinWeight(newCheckin, weight)
}

// 将更新后打卡人数添加到zset
func ZaddCheckinJoinNum(cid int, checkin *model.Checkin) (err error) {
	return redis.ZaddCheckinJoinNum(cid, checkin)
}

// 获取zset缓存
func GetZsetCheckinNum() (zrankm map[string]int, err error) {
	return redis.GetZsetCheckinNum()
}
