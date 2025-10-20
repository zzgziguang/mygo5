package redis

import (
	"context"
	"demo1/internal/app/mydemo/model"
	"encoding/json"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/mitchellh/mapstructure"
)

var Ctx = context.Background()
var ttl = 5 * time.Minute

// 从 Redis 缓存获取用户
func GetUserFromCache(id int) (user *model.User, err error) {
	key := "user:" + strconv.Itoa(id)
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
	user = &model.User{}
	// 使用 mapstructure 将 map 转为结构体（支持类型转换）
	if err = mapstructure.WeakDecode(data, &user); err != nil {
		return
	}
	return
}

// 保存用户到 Redis 缓存
func SetUserToCache(user *model.User) error {
	key := "user:" + strconv.Itoa(user.Id)
	var data map[string]interface{}
	if err := mapstructure.WeakDecode(user, &data); err != nil {
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
func DelRedisUser(idStr string) (err error) {

	cacheKey := "user:" + idStr
	err = RedisClient.Del(Ctx, cacheKey).Err()
	return
}

// 查询列表缓存
func GetRedisUserSlice(order string, page int, pagesize int) (users []model.User, err error) {
	sliceKey := "users:" + order + strconv.Itoa(page) + strconv.Itoa(pagesize)
	strSlice, err := RedisClient.Get(Ctx, sliceKey).Result()
	if err != nil {
		if err == redis.Nil {
			err = nil
			return
		}
	}
	date := []byte(strSlice)
	err = json.Unmarshal(date, &users)
	if err != nil {
		return
	}
	return
}

// 添加列表缓存
func SetRedisUserSlice(users []model.User, order string, page int, pagesize int) (err error) {
	if users == nil { //users不能等于nil
		return
	}
	key := "users:" + order + strconv.Itoa(page) + strconv.Itoa(pagesize)
	byteSlice, err := json.Marshal(users)
	if err != nil {
		return
	}
	str := string(byteSlice)
	err = RedisClient.Set(Ctx, key, str, ttl).Err()
	if err != nil {
		return
	}
	return
}

// 从缓存获取用户数量
func GetRedisUserCount() (total int64, err error) {
	key := "usercount"
	count, err := RedisClient.Get(Ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			err = nil
			return
		} else {
			return
		}

	}
	num, err := strconv.Atoi(count)
	if err != nil {
		return
	}
	total = int64(num)
	return
}

// 添加用户数量到缓存
func SetRedisUserCount(total int64) (err error) {
	key := "usercount"
	err = RedisClient.Set(Ctx, key, total, ttl).Err()
	if err != nil {
		return
	}
	return
}
