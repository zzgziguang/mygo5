package service

import (
	"demo1/internal/app/mydemo/dao/mysql"
	"demo1/internal/app/mydemo/dao/redis"
	"demo1/internal/app/mydemo/model"
	"time"
)

// 查询表
func GetUserCheckinJoin(uid int, cid int) (userCheckinJoin *model.UserCheckinJoin, err error) {
	return mysql.GetUserCheckinJoin(uid, cid)
}

// 根据uid查询表
func GetUserCheckinJoinByuid(uid int) (userCheckinJoinByuid []*model.UserCheckinJoin, err error) {
	return mysql.GetUserCheckinJoinByuid(uid)
}

// 查询表
func GetUserCheckinJoinCount(uid int, cid int) (userCheckinJoin *model.UserCheckinJoin, err error) {
	return mysql.GetUserCheckinJoinCount(uid, cid)
}

// 添加表
func AddUserCheckinJoin(newUserCheckinJoin *model.UserCheckinJoin) (err error) {
	return mysql.AddUserCheckinJoin(newUserCheckinJoin)
}

// 从 Redis 缓存获取用户
func GetUserCheckinJoinFromCache(uid int, cid int) (userRedisCheckinJoin *model.UserCheckinJoin, err error) {
	return redis.GetUserCheckinJoinFromCache(uid, cid)
}

// 保存用户到 Redis 缓存
func SetUserCheckinJoinToCache(userRedisCheckinJoin *model.UserCheckinJoin, ttl time.Duration) error {
	return redis.SetUserCheckinJoinToCache(userRedisCheckinJoin, ttl)
}
