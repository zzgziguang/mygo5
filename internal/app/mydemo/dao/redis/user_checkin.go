package redis

import (
	"demo1/internal/app/mydemo/model"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/mitchellh/mapstructure"
)

// 添加usercheckin hset
func HSetUserCheckinToCache(uid int, userCheckin *model.UserCheckin) (err error) {
	key := "uid:" + strconv.Itoa(uid)
	var data map[string]interface{}

	err = mapstructure.WeakDecode(userCheckin, &data)
	if err != nil {
		return err
	}

	err = RedisClient.HSet(Ctx, key, data).Err()
	if err != nil {
		return
	}
	return
}

// 获取usercheckin
func HGetUserCheckinFromCache(uid int) (userCheckin *model.UserCheckin, err error) {
	key := "uid:" + strconv.Itoa(uid)
	data, err := RedisClient.HGetAll(Ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			err = nil
			return
		}
		return
	}

	if len(data) == 0 {
		return
	}

	checkinLastTimeMap := make(map[int64]time.Time, 0)
	checkinJoinTimeMap := make(map[int64]time.Time, 0)
	checkinDateNumTimeMap := make(map[int64]int, 0)
	var dayNum int
	var lastDate string
	for k, v := range data {
		ok := strings.HasPrefix(k, "lasttime")
		if ok {
			cidStr := strings.TrimLeft(k, "lasttime")
			var cid int64
			cid, err = strconv.ParseInt(cidStr, 10, 64)
			if err != nil {
				return
			}

			var vInt int64
			vInt, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				return
			}

			checkinLastTimeMap[cid] = time.Unix(vInt, 0)
		}

		ok = strings.HasPrefix(k, "jointime")
		if ok {
			cidStr := strings.TrimLeft(k, "jointime")
			var cid int64
			cid, err = strconv.ParseInt(cidStr, 10, 64)
			if err != nil {
				return
			}

			var vInt int64
			vInt, err = strconv.ParseInt(v, 10, 64)
			if err != nil {
				return
			}

			checkinJoinTimeMap[cid] = time.Unix(vInt, 0)
		}

		ok = strings.HasPrefix(k, "datenum")
		if ok {
			cidStr := strings.TrimLeft(k, "datenum")
			var cid int64
			cid, err = strconv.ParseInt(cidStr, 10, 64)
			if err != nil {
				return
			}

			var vInt int
			vInt, err = strconv.Atoi(v)
			if err != nil {
				return
			}

			checkinDateNumTimeMap[cid] = vInt
		}
		if k == "num" {
			var vInt int
			vInt, err = strconv.Atoi(v)
			if err != nil {
				return
			}
			dayNum = vInt
		}

		if k == "lastdate" {
			lastDate = v
		}
	}
	userCheckin = &model.UserCheckin{}

	// var dayNum int = 0
	// for k, v := range checkinDateNumTimeMap {
	// 	dayNum = dayNum + v
	// }

	userCheckin.CheckinDayNumMap = checkinDateNumTimeMap
	userCheckin.CheckinJoinTimeMap = checkinJoinTimeMap
	userCheckin.CheckinLastTimeMap = checkinLastTimeMap
	userCheckin.DayNum = dayNum
	userCheckin.LastDate = lastDate

	return
}

// 使用字段添加usercheckin hset
func HSetFieldUserCheckinToCache(uid int, cid int, fieldName string, value interface{}) (err error) {
	key := "uid:" + strconv.Itoa(uid)
	field := fieldName + strconv.Itoa(cid)

	err = RedisClient.HSet(Ctx, key, field, value).Err()
	if err != nil {
		return
	}
	return
}

// 添加usercheckin jointime hset
func HSetUserCheckinJoinTimeToCache(uid int, cid int, joinTime int64) (err error) {
	key := "uid:" + strconv.Itoa(uid)
	field := "jointime" + strconv.Itoa(cid)

	err = RedisClient.HSetNX(Ctx, key, field, joinTime).Err()
	if err != nil {
		return
	}
	return
}

// 添加usercheckin lasttime hset
func HSetUserCheckinLastTimeToCache(uid int, cid int, lastTime int64) (err error) {
	key := "uid:" + strconv.Itoa(uid)
	field := "lasttime" + strconv.Itoa(cid)

	err = RedisClient.HSet(Ctx, key, field, lastTime).Err()
	if err != nil {
		return
	}
	return
}

// 添加usercheckin datenum hset
func HSetUserCheckinDateNumToCache(uid int, cid int) (err error) {
	key := "uid:" + strconv.Itoa(uid)
	field := "datenum" + strconv.Itoa(cid)

	err = RedisClient.HIncrBy(Ctx, key, field, 1).Err()
	if err != nil {
		return
	}
	return
}

// 字段lastdate添加
func HSetUserCheckinLastDateToCache(uid int, dateTime string) (err error) {
	key := "uid:" + strconv.Itoa(uid)
	field := "lastdate"

	err = RedisClient.HSet(Ctx, key, field, dateTime).Err()
	if err != nil {
		return
	}
	return
}

// daynum添加
func HSetUserCheckinDayNumToCache(uid int) (err error) {
	key := "uid:" + strconv.Itoa(uid)
	field := "num"

	err = RedisClient.HIncrBy(Ctx, key, field, 1).Err()
	if err != nil {
		return
	}
	return
}
