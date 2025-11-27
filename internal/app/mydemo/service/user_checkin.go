package service

import (
	"demo1/internal/app/mydemo/dao/redis"
	"demo1/internal/app/mydemo/model"
)

// 获取uid缓存
func HGetUserCheckinFromCache(uid int) (*model.UserCheckin, error) {
	return redis.HGetUserCheckinFromCache(uid)
}

// 添加usercheckin jointime hset某打卡参与打卡时间
func HSetUserCheckinJoinTimeToCache(uid int, cid int, joinTime int64) (err error) {
	return redis.HSetUserCheckinJoinTimeToCache(uid, cid, joinTime)
}

// 添加usercheckin lasttime hset某打卡上次打卡时间
func HSetUserCheckinLastTimeToCache(uid int, cid int, lastTime int64) (err error) {
	return redis.HSetUserCheckinLastTimeToCache(uid, cid, lastTime)
}

// 添加usercheckin datenum hset某打卡打卡天数
func HSetUserCheckinDateNumToCache(uid int, cid int) (err error) {
	return redis.HSetUserCheckinDateNumToCache(uid, cid)
}

// 字段lastdate添加 用户上次打卡日期
func HSetUserCheckinLastDateToCache(uid int, dateTime string) (err error) {
	return redis.HSetUserCheckinLastDateToCache(uid, dateTime)
}

// daynum添加 用户打卡总天数
func HSetUserCheckinDayNumToCache(uid int) (err error) {
	return redis.HSetUserCheckinDayNumToCache(uid)
}
