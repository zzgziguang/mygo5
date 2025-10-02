package mysql

import (
	"demo1/internal/app/mydemo/model"

	"gorm.io/gorm"
)

func CreateCheckin(newcheckin *model.Checkin) (result *gorm.DB) {
	result = DB.Create(&newcheckin)
	return
}
