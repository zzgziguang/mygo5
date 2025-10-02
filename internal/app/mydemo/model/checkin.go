package model

import "time"

type Checkin struct {
	Id            int       `json:"id" gorm:"column:id" mapstructure:"id"`
	Title         string    `json:"title" gorm:"column:title" mapstructure:"title"`
	CreateAt      time.Time `json:"createAt" gorm:"column:create_at" mapstructure:"createAt"`
	UpdateAt      time.Time `json:"updateAt" gorm:"column:update_at" mapstructure:"updateAt"`
	CheckinStatus int       `json:"checkinStatus" gorm:"column:checkin_status" mapstructure:"checkinStatus"`
}

const (
	CheckinDelete int = 0 //删除
	CheckinNormal int = 1 //正常
	CheckinReview int = 2 //审核
)

// 指定Checkin对应的表名
func (Checkin) TableName() string {
	return "checkin"
}
