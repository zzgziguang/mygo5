package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// 初始化数据库连接
func DaoInitDB(dsn string) (err error) {
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	return
}
