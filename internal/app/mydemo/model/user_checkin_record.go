package model

import "time"

type UserCheckinRecord struct {
	Id       int        `json:"id" gorm:"column:id" mapstructure:"id"`
	Uid      int        `json:"uid" gorm:"column:uid" mapstructure:"uid"`
	Cid      int        `json:"cid" gorm:"column:cid" mapstructure:"cid"`
	Date     int        `json:"date" gorm:"column:date" mapstructure:"date"`
	CreateAt *time.Time `json:"createAt" gorm:"column:create_at" mapstructure:"createAt"`
	UpdateAt *time.Time `json:"updateAt" gorm:"column:update_at" mapstructure:"updateAt"`
	Status   int        `json:"status" gorm:"column:status" mapstructure:"status"`
}

const (
	RecordStatusDelete int = 0 //删除
	RecordStatusNorMal int = 1 //正常
)

func (UserCheckinRecord) TableName() string {
	return "user_checkin_record"
}
