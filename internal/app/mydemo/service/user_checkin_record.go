package service

import (
	"demo1/internal/app/mydemo/dao/mysql"
	"demo1/internal/app/mydemo/dao/redis"
	"demo1/internal/app/mydemo/model"
	"time"
)

// 查询有没有打卡
func GetUserCheckinRecord(uid int, cid int, date int) (userCheckinRecord *model.UserCheckinRecord, err error) {
	return mysql.GetUserCheckinRecord(uid, cid, date)
}

// 根据uid查询
func GetUserCheckinRecordByUidDate(uid int, date int) (userCheckinRecordByuid []model.UserCheckinRecord, err error) {
	return mysql.GetUserCheckinRecordByUidDate(uid, date)
}

// in 查询
func GetUserCheckinRecordInCheckinId(uid int, date int, cidSlice []int) (userCheckinRecordByuid []model.UserCheckinRecord, err error) {
	return mysql.GetUserCheckinRecordInCheckinId(uid, date, cidSlice)
}

// 添加打卡
func AddUserCheckinRecord(userCheckinRecord *model.UserCheckinRecord) (err error) {
	return mysql.AddUserCheckinRecord(userCheckinRecord)
}

// 打卡列表
func GetUserCheckinRecordList(uid int, cid int, isasc bool) (userCheckinRecordList []model.UserCheckinRecord, err error) {
	return mysql.GetUserCheckinRecordList(uid, cid, isasc)
}

// 更新打卡
func UpdateUserCheckinRecord(uid int, cid int, date int) (err error) {
	return mysql.UpdateUserCheckinRecord(uid, cid, date)
}

// 从 Redis 缓存获取用户
func GetUserCheckinRecordFromCache(uid int, cid int, date int) (userRedisCheckinRecord *model.UserCheckinRecord, err error) {
	return redis.GetUserCheckinRecordFromCache(uid, cid, date)
}

// 保存用户到 Redis 缓存
func SetUserCheckinRecordToCache(userRedisCheckinRecord *model.UserCheckinRecord, ttl time.Duration) error {
	return redis.SetUserCheckinRecordToCache(userRedisCheckinRecord, ttl)
}

// 更新后删除 Redis 缓存，下次查询会重建
func DelRedisUserCheckinRecord(uid int, cid int, date int) (err error) {
	return redis.DelRedisUserCheckinRecord(uid, cid, date)
}
