package service

import (
	"demo1/internal/app/mydemo/dao/mysql"
	"demo1/internal/app/mydemo/dao/redis"
	"demo1/internal/app/mydemo/model"

	"gorm.io/gorm"
)

// 从 Redis 缓存获取用户
func GetUserFromCache(id int) (*model.User, error) {
	return redis.GetUserFromCache(id)

}

// 保存用户到 Redis 缓存
func SetUserToCache(user *model.User) error {
	return redis.SetUserToCache(user)
}

// 更新后删除 Redis 缓存，下次查询会重建
func DelRedisUser(idStr string) (err error) {
	return redis.DelRedisUser(idStr)
}

// 从缓存获取用户列表
func GetRedisUserSlice(order string, page int, pagesize int) (users []model.User, err error) {
	return redis.GetRedisUserSlice(order, page, pagesize)
}

// 保存用户列表到缓存
func SetRedisUserSlice(users []model.User, order string, page int, pagesize int) (err error) {
	return redis.SetRedisUserSlice(users, order, page, pagesize)
}

// 从缓存获取用户数量
func GetRedisUserCount() (total int64, err error) {
	return redis.GetRedisUserCount()
}

// 添加用户数量到缓存
func SetRedisUserCount(total int64) (err error) {
	return redis.SetRedisUserCount(total)
}

// 在数据库添加用户
func CreateUser(newUser *model.User) (result *gorm.DB) {
	return mysql.CreateUser(newUser)
}

// 从mysql根据id查询用户
func GerUserById(id int) (result *gorm.DB, dbUser *model.User) {
	return mysql.GerUserById(id)
}

// 在mysql更新数据
func UpdateUserById(id int, username string) (result *gorm.DB) {
	return mysql.UpdateUserById(id, username)
}

// 在mysql查询并排序
func GetUserByPage(page int, pagesize int, isasc bool) (users []model.User, err error) {
	return mysql.GetUserByPage(page, pagesize, isasc)
}

// 查询mysql数据个数
func GetUserCount() (total int64, err error) {
	return mysql.GetUserCount()
}
