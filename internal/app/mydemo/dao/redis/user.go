package redis

import (
	"context"
	"demo1/internal/app/mydemo/model"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/mitchellh/mapstructure"
)

var Ctx = context.Background()

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
	user = &model.User{}
	// 使用 mapstructure 将 map 转为结构体（支持类型转换）
	if err = mapstructure.WeakDecode(data, &user); err != nil {
		return nil, err
	}
	return user, nil
}

// 保存用户到 Redis 缓存
func SetUserToCache(user *model.User, ttl time.Duration) error {
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
