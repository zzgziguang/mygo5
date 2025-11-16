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
func SetUserCheckinRecordToCache(userRedisCheckinRecord *model.UserCheckinRecord, ttl time.Duration) (err error) {
	key := "checkinRecord:" + strconv.Itoa(userRedisCheckinRecord.Uid) + strconv.Itoa(userRedisCheckinRecord.Cid)
	var data map[string]interface{}
	if err := mapstructure.WeakDecode(userRedisCheckinRecord, &data); err != nil {
		return err
	}
	// 写入 hash
	err = RedisClient.HSet(Ctx, key, data).Err()
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

// 添加名次到zset缓存
func ZaddUserCheckinRecordCountToCache(cid int, uid int, recordTime time.Time, date int) (err error) {
	key := "checkinDate" + strconv.Itoa(cid) + strconv.Itoa(date)
	members := redis.Z{
		Score:  float64(recordTime.Unix()),
		Member: uid,
	}

	err = RedisClient.ZAdd(Ctx, key, &members).Err()
	return
}

// 获取打卡zset缓存名次
func ZrankUserCheckinRecordCountToCache(cid int, uid int, date int) (rank int, err error) {
	key := "checkinDate" + strconv.Itoa(cid) + strconv.Itoa(date)
	member := strconv.Itoa(uid)

	zrank, err := RedisClient.ZRank(Ctx, key, member).Result()
	if err != nil {
		return
	}

	rank = int(zrank) + 1
	return
}

// 添加名次到缓存string
func IncrUserCheckinRecordCountToCache(cid int, date int) (rank int64, err error) {
	key := "checkinDateFromString" + strconv.Itoa(cid) + strconv.Itoa(date)
	rank, err = RedisClient.Incr(Ctx, key).Result()
	return
}
