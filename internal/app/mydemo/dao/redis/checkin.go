package redis

import (
	"demo1/internal/app/mydemo/model"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/mitchellh/mapstructure"
)

// 获取打卡
func GetCheckinFromCache(cid int) (checkin *model.Checkin, err error) {
	key := "checkin:" + strconv.Itoa(cid)
	data, err := RedisClient.HGetAll(Ctx, key).Result()
	if err != nil {
		return
	}

	checkin = &model.Checkin{}

	// 使用 mapstructure 将 map 转为结构体
	err = mapstructure.WeakDecode(data, checkin)
	return

}

// 保存checkin到 Redis 缓存
func SetCheckinToCache(checkin *model.Checkin, ttl time.Duration) (err error) {
	key := "checkin:" + strconv.Itoa(checkin.Id)
	var data map[string]interface{}

	err = mapstructure.WeakDecode(checkin, &data)
	if err != nil {
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

// 将更新后打卡人数添加到zset
func ZaddCheckinJoinNum(cid int, checkin *model.Checkin) (err error) {
	key := "checkinNumber"
	members := redis.Z{
		Score:  float64(checkin.JoinNum),
		Member: cid,
	}

	err = RedisClient.ZAdd(Ctx, key, &members).Err()
	return
}

// 获取全部zset缓存排名
func GetZsetCheckinNum() (zrankm map[string]int, err error) {
	key := "checkinNumber"
	var start int64 = 0
	var stop int64 = -1
	zrankm = make(map[string]int, 0)

	zranks, err := RedisClient.ZRange(Ctx, key, start, stop).Result()
	if err != nil {
		if err == redis.Nil {
			err = nil
			return
		}
	}

	for i, v := range zranks {
		zrankm[v] = i + 1
	}
	return
}

// 更新参与打卡的权重
func ZaddCheckinWeight(cid int, checkin *model.Checkin) (err error) {
	key := "checkinWeight"
	members := redis.Z{
		Score:  float64(checkin.Weight),
		Member: cid,
	}

	err = RedisClient.ZAdd(Ctx, key, &members).Err()
	return
}

// 获取全部zset缓存排名
func GetZsetCheckinWeight() (zrankm map[string]int, err error) {
	key := "checkinWeight"
	var start int64 = 0
	var stop int64 = -1
	zrankm = make(map[string]int, 0)

	zranks, err := RedisClient.ZRange(Ctx, key, start, stop).Result()
	if err != nil {
		if err == redis.Nil {
			err = nil
			return
		}
	}

	for i, v := range zranks {
		zrankm[v] = i + 1
	}
	return
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
