package service

import (
	"demo1/internal/app/mydemo/dao/kafka"
	"demo1/internal/app/mydemo/dao/mysql"
	"demo1/internal/app/mydemo/dao/redis"
	"demo1/internal/app/mydemo/model"
)

// 查询表
func GetUserCheckinJoin(uid int, cid int) (userCheckinJoin *model.UserCheckinJoin, err error) {
	return mysql.GetUserCheckinJoin(uid, cid)
}

// 根据uid查询表
func GetUserCheckinJoinByuid(uid int) (userCheckinJoinByuid []model.UserCheckinJoin, err error) {
	return mysql.GetUserCheckinJoinByuid(uid)
}

// 获取参与打卡用户uid
func GetCreateUidFromUserCheckinJoin() (uids []int, err error) {
	return mysql.GetCreateUidFromUserCheckinJoin()
}

// 查询表
func GetUserCheckinJoinCount(uid int, cid int) (userCheckinJoin *model.UserCheckinJoin, err error) {
	return mysql.GetUserCheckinJoinCount(uid, cid)
}

// 获取参与打卡用户uid,cid
func GetCreateUidCidFromUserCheckinJoin(page int, pagesize int, isasc bool) (userCheckinJoins []model.UserCheckinJoinUidCid, err error) {
	return mysql.GetCreateUidCidFromUserCheckinJoin(page, pagesize, isasc)
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
func SetUserCheckinJoinToCache(userRedisCheckinJoin *model.UserCheckinJoin) (err error) {
	return redis.SetUserCheckinJoinToCache(userRedisCheckinJoin)
}

// 添加kafka的数据
func ProducerSend(value []byte) (partition int32, offset int64, err error) {
	return kafka.ProducerSend(value)
}
