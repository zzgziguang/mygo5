package model

import "time"

type UserCheckinJoin struct {
	Id       int       `json:"id" gorm:"column:id" mapstructure:"id"`
	Uid      int       `json:"uid" gorm:"column:uid" mapstructure:"uid"`
	Cid      int       `json:"cid" gorm:"column:cid" mapstructure:"cid"`
	JoinTime *time.Time `json:"joinTime" gorm:"column:joinTime" mapstructure:"joinTime"`
	CreateAt *time.Time `json:"createAt" gorm:"column:createAt" mapstructure:"createAt"`
	UpdateAt *time.Time `json:"updateAt" gorm:"column:updateAt" mapstructure:"updateAt"`
	Status   int       `json:"status" gorm:"column:status" mapstructure:"status"`
}

const (
	JoinStatusDelete int = 0 //删除
	JoinStatusNormal int = 1 //正常
)
func (UserCheckinJoin)TableName string{
	return "user_checkin_join"
}
