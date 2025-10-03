package redis

import (
	"demo1/internal/app/mydemo/model"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
)

//var Ctx = context.Background()

// 保存checkin到 Redis 缓存
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

/*
// 查询checkin在Redis缓存数据
func GetCheckinFormCache(id int) (checkin *model.Checkin, err error) {
	key := strconv.Itoa(id)
	res, err := RedisClient.HGetAll(Ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			err = nil
			return
		} else {
			return
		}
	}
	checkin = &model.Checkin{} //给checkin空间
	err = mapstructure.WeakDecode(res, &checkin)
	if err != nil {
		return
	}
	return
}
*/
