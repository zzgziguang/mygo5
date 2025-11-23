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

// 查询最近两天打卡条数
func GetUserCheckinRecordTwoDayCount(uid int, cid int, yestoday int) (count int64, err error) {
	return mysql.GetUserCheckinRecordTwoDayCount(uid, cid, yestoday)
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

// 获取今天的用户打卡的时间
func GetCreateTimeFromUserCheckinRecord(uid int, cid int, date int) (createTime time.Time, err error) {
	return mysql.GetCreateTimeFromUserCheckinRecord(uid, cid, date)
}

// 获取用户该打卡在今天打卡的名次
func GetUserCheckinRecordByCount(cid int, date int, createTime time.Time) (rank int, err error) {
	return mysql.GetUserCheckinRecordByCount(cid, date, createTime)
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
func SetUserCheckinRecordToCache(userRedisCheckinRecord *model.UserCheckinRecord, ttl time.Duration) (err error) {
	return redis.SetUserCheckinRecordToCache(userRedisCheckinRecord, ttl)
}

// 更新后删除 Redis 缓存，下次查询会重建
func DelRedisUserCheckinRecord(uid int, cid int, date int) (err error) {
	return redis.DelRedisUserCheckinRecord(uid, cid, date)
}

// 添加名次到zset缓存
func ZaddUserCheckinRecordCountToCache(cid int, uid int, recordTime time.Time, date int) (err error) {
	return redis.ZaddUserCheckinRecordCountToCache(cid, uid, recordTime, date)
}

// 获取打卡zset缓存名次
func ZrankUserCheckinRecordCountToCache(cid int, uid int, date int) (rank int, err error) {
	return redis.ZrankUserCheckinRecordCountToCache(cid, uid, date)
}

// 添加名次到缓存string
func IncrUserCheckinRecordCountToCache(cid int, date int) (rank int64, err error) {
	return redis.IncrUserCheckinRecordCountToCache(cid, date)
}
