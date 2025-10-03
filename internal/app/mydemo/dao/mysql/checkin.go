package mysql

import (
	"demo1/internal/app/mydemo/model"

	"gorm.io/gorm"
)

func CreateCheckin(newcheckin *model.Checkin) (result *gorm.DB) {
	result = DB.Create(&newcheckin)
	return
}
func GetCheckinOrderId(page int, pagesize int, isasc bool) (checkins []model.Checkin, err error) {
	offset := (page - 1) * pagesize
	var order string
	if isasc {
		order = "id asc"
	} else {
		order = "id desc"
	}
	err = DB.Order(order).Offset(offset).Limit(pagesize).Find(&checkins).Error
	return
}
