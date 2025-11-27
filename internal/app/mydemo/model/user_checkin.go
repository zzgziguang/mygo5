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

// 1、是否今日首次打卡
// 2、某打卡是否首次打卡
// 3、某打卡是否今日首次打卡

//是否今日首次打卡
func (u UserCheckin) TodayFristCheckinBool() (bool bool) {
	daytime := time.Now()
	if u.LastDate != daytime.Format("20060102") {
		bool = true
	} else {
		bool = false
	}

	return
}

//某打卡是否首次打卡
func (u UserCheckin) FristJoinCheckinBool(cid int64) (bool bool) {
	dayNum, ok := u.CheckinDayNumMap[cid]
	if ok {
		if dayNum == 0 {
			bool = false
		} else {
			bool = true
		}
	}
	return
}

//某打卡是否今日首次打卡
func (u UserCheckin) FristDateJoinCheckinBool(cid int64) (bool bool) {
	lastTime, ok := u.CheckinLastTimeMap[cid]

	if ok {

		if lastTime.YearDay() < time.Now().YearDay() {
			bool = true
		} else {
			bool = false
		}
	}
	return
}
