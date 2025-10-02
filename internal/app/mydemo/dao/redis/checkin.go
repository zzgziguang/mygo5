package redis

import (
	"demo1/internal/app/mydemo/model"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
)

// 保存用户到 Redis 缓存
func SetCheckinToCache(checkin *model.Checkin, ttl time.Duration) error {
	key := "user:" + strconv.Itoa(checkin.Id)
	var data map[string]interface{}
	if err := mapstructure.WeakDecode(checkin, &data); err != nil {
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
