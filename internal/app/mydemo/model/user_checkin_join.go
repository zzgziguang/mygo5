package model

import "time"

type UserCheckinJoin struct {
	Id          int        `json:"id" gorm:"column:id" mapstructure:"id"`
	Uid         int        `json:"uid" gorm:"column:uid" mapstructure:"uid"`
	Cid         int        `json:"cid" gorm:"column:cid" mapstructure:"cid"`
	JoinTime    *time.Time `json:"joinTime" gorm:"column:join_time" mapstructure:"-"`
	JoinTimeStr string     `json:"-" gorm:"-" mapstructure:"joinTime"`
	CreateAt    *time.Time `json:"createAt" gorm:"column:create_at" mapstructure:"-"` //mapstructure:"-"
	CreateAtStr string     `json:"-" gorm:"-" mapstructure:"createAt"`
	UpdateAt    *time.Time `json:"updateAt" gorm:"column:update_at" mapstructure:"-"`
	UpdateAtStr string     `json:"-" gorm:"-" mapstructure:"updateAt"` //json"-"表示去掉这个字段
	Status      int        `json:"status" gorm:"column:status" mapstructure:"status"`
}

type UserCheckinJoinUidCid struct {
	Uid int `gorm:"column:uid"`
	Cid int `gorm:"column:cid"`
}

const (
	JoinStatusDelete int = 0 //删除
	JoinStatusNormal int = 1 //正常
)

func (UserCheckinJoin) TableName() string {
	return "user_checkin_join"
}
