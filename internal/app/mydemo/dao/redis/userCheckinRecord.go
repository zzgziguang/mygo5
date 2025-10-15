package redis

import (
	"demo1/internal/app/mydemo/model"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/mitchellh/mapstructure"
)

// 从 Redis 缓存获取用户
func GetUserCheckinRecordFromCache(uid int, cid int, date int) (userRedisCheckinRecord *model.UserCheckinRecord, err error) {
	key := "checkinRecord:" + strconv.Itoa(uid) + strconv.Itoa(cid) + strconv.Itoa(date)
	data, err := RedisClient.HGetAll(Ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			err = nil
			return nil, nil
		} else {
			return nil, err
		}
	} else {
		if len(data) > 0 {
			userRedisCheckinRecord = &model.UserCheckinRecord{}
			err = mapstructure.WeakDecode(data, &userRedisCheckinRecord)
			if err != nil {
				return nil, err
			} else {
				return userRedisCheckinRecord, nil
			}
		} else {
			return nil, nil
		}
	}
}

// 保存用户到 Redis 缓存
func SetUserCheckinRecordToCache(userRedisCheckinRecord *model.UserCheckinRecord, ttl time.Duration) error {
	key := "checkinRecord:" + strconv.Itoa(userRedisCheckinRecord.Uid) + strconv.Itoa(userRedisCheckinRecord.Cid)
	var data map[string]interface{}
	if err := mapstructure.WeakDecode(userRedisCheckinRecord, &data); err != nil {
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

// 更新后删除 Redis 缓存，下次查询会重建
func DelRedisUserCheckinRecord(uid int, cid int, date int) (err error) {

	cacheKey := "checkinRecord" + strconv.Itoa(uid) + strconv.Itoa(cid)
	err = RedisClient.Del(Ctx, cacheKey).Err()
	return
}
