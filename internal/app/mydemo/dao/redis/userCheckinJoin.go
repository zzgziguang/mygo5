package redis

import (
	"demo1/internal/app/mydemo/model"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/mitchellh/mapstructure"
)

// 从 Redis 缓存获取用户
func GetUserCheckinJoinFromCache(uid int, cid int) (userRedisCheckinJoin *model.UserCheckinJoin, err error) {
	key := "checkinJoin:" + strconv.Itoa(uid) + strconv.Itoa(cid)
	data, err := RedisClient.HGetAll(Ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			err = nil
			return
		} else {
			return
		}
	}

	if len(data) == 0 {
		return
	}

	userRedisCheckinJoin = &model.UserCheckinJoin{}
	// 使用 mapstructure 将 map 转为结构体
	err = mapstructure.WeakDecode(data, userRedisCheckinJoin)
	if err != nil {
		return
	}
	//将hash中jointimestr字符串转为时间格式jointime
	userRedisCheckinJoinTime, err := time.ParseInLocation("2006-01-02 15:04:05", userRedisCheckinJoin.JoinTimeStr, time.Local)
	if err != nil {
		return
	} else {
		userRedisCheckinJoin.JoinTime = &userRedisCheckinJoinTime
	}
	//将hash中createatstr字符串转为时间格式createat
	userRedisCheckinJoinCreateAt, err := time.ParseInLocation("2006-01-02 15:04:05", userRedisCheckinJoin.CreateAtStr, time.Local)
	if err != nil {
		return
	} else {
		userRedisCheckinJoin.CreateAt = &userRedisCheckinJoinCreateAt
	}
	//将hash中updateatstr字符串转为时间格式updateat
	userRedisCheckinJoinUpdateAt, err := time.ParseInLocation("2006-01-02 15:04:05", userRedisCheckinJoin.UpdateAtStr, time.Local)
	if err != nil {
		return
	} else {
		userRedisCheckinJoin.UpdateAt = &userRedisCheckinJoinUpdateAt
	}

	return
}

// 保存用户到 Redis 缓存
func SetUserCheckinJoinToCache(userRedisCheckinJoin *model.UserCheckinJoin, ttl time.Duration) error {
	key := "checkinJoin:" + strconv.Itoa(userRedisCheckinJoin.Uid) + strconv.Itoa(userRedisCheckinJoin.Cid)
	var data map[string]interface{}
	//将结构体中时间转为字符串
	userRedisCheckinJoin.JoinTimeStr = userRedisCheckinJoin.JoinTime.Format("2006-01-02 15:04:05")
	userRedisCheckinJoin.CreateAtStr = userRedisCheckinJoin.CreateAt.Format("2006-01-02 15:04:05")
	userRedisCheckinJoin.UpdateAtStr = userRedisCheckinJoin.UpdateAt.Format("2006-01-02 15:04:05")
	if err := mapstructure.WeakDecode(userRedisCheckinJoin, &data); err != nil {
		return err
	}

	// 写入 hash
	err := RedisClient.HSet(Ctx, key, data).Err()
	if err != nil {
		return err
	}
	RedisClient.Expire(Ctx, key, ttl)
	return nil
}
