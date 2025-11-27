package service

import (
	"demo1/internal/app/mydemo/dao/redis"
	"demo1/internal/app/mydemo/model"
)

func HGetUserCheckinFromCache(uid int) (*model.UserCheckin, error) {
	return redis.HGetUserCheckinFromCache(uid)
}

// 添加usercheckin jointime hset
func HSetUserCheckinJoinTimeToCache(uid int, cid int, joinTime int64) (err error) {
	return redis.HSetUserCheckinJoinTimeToCache(uid, cid, joinTime)
}

// 添加usercheckin lasttime hset
func HSetUserCheckinLastTimeToCache(uid int, cid int, lastTime int64) (err error) {
	return redis.HSetUserCheckinLastTimeToCache(uid, cid, lastTime)
}

// 添加usercheckin datenum hset
func HSetUserCheckinDateNumToCache(uid int, cid int) (err error) {
	return redis.HSetUserCheckinDateNumToCache(uid, cid)
}

// 字段lastdate添加
func HSetUserCheckinLastDateToCache(uid int, dateTime string) (err error) {
	return redis.HSetUserCheckinLastDateToCache(uid, dateTime)
}

// daynum添加
func HSetUserCheckinDayNumToCache(uid int) (err error) {
	return redis.HSetUserCheckinDayNumToCache(uid)
}
