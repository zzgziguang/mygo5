package mysql

import (
	"demo1/internal/app/mydemo/model"

	"gorm.io/gorm"
)

func CreateUser(newUser model.User) (result *gorm.DB) {
	result = DB.Create(&newUser)
	return
}
func GerUserById(id int) (result *gorm.DB, dbUser model.User) {
	result = DB.First(&dbUser, id)
	return
}
func UpdateUserById(id int, username string) (result *gorm.DB) {
	result = DB.Model(&model.User{}).Where("id = ?", id).Update("username", username)
	return
}
func GetUserByPage(page int, pagesize int, isasc bool) (users []model.User, err error) {
	offset := (page - 1) * pagesize
	orderby := "id desc"
	if isasc {
		orderby = "id asc"
	}
	err = DB.Order(orderby).Offset(offset).Limit(pagesize).Find(&users).Error
	return
}
