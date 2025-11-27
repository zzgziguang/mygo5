package model

import "time"

type UserCheckin struct {
	//上次打卡日期
	LastDate string
	//总打卡天数
	DayNum int
	//每个打卡的天数 cid => num
	CheckinDayNumMap map[int64]int
	//每个打卡的参与时间 cid => joinTime
	CheckinJoinTimeMap map[int64]time.Time
	//每个打卡的上次打卡时间 cid => lastTime
	CheckinLastTimeMap map[int64]time.Time
}
