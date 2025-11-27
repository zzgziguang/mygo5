package model

import "time"

type Checkin struct {
	Id    int    `json:"id" gorm:"column:id" mapstructure:"id"`
	Title string `json:"title" gorm:"column:title" mapstructure:"title"`
	// CreateAt      *time.Time `json:"createAt" gorm:"column:create_at" mapstructure:"createAt"`
	CreateAt      *time.Time `json:"createAt" gorm:"column:create_at" mapstructure:"-"`
	CreateAtStr   string     `json:"-" gorm:"-" mapstructure:"createAt"`
	UpdateAt      *time.Time `json:"updateAt" gorm:"column:update_at" mapstructure:"-"`
	UpdateAtStr   string     `json:"-" gorm:"-" mapstructure:"updateAt"`
	CheckinStatus int        `json:"checkinStatus" gorm:"column:checkin_status" mapstructure:"checkinStatus"`
	JoinNum       int        `json:"join_num" gorm:"column:join_num" mapstructure:"join_num"`
	Weight        int        `json:"weight" gorm:"column:weight" mapstructure:"weight"`
	StartTime     int64      `json:"startTime" gorm:"column:start_time" mapstructure:"startTime"`
	EndTime       int64      `json:"endTime" gorm:"column:end_time" mapstructure:"endTime"`
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

//打卡排序结构体
type CheckinSort struct {
	Checkin
	JoinBool   bool      `json:"-"`
	RecordBool bool      `json:"-"`
	JoinTime   time.Time `json:"-"`
	Rank       int       `json:"rank"`
}

// 排序
type CheckinSlice []CheckinSort

//  < 升
func (c CheckinSlice) Len() int { return len(c) }

func (c CheckinSlice) Less(i, j int) bool {
	// 参与打卡的
	if c[i].JoinBool && !c[j].JoinBool {
		return true
	} else if !c[i].JoinBool && c[j].JoinBool {
		return false
	}

	//全部参与
	if c[i].JoinBool && c[j].JoinBool {
		if c[i].RecordBool && !c[j].RecordBool {
			return true
		} else if !c[i].RecordBool && c[j].RecordBool {
			return false
		} else {
			return c[i].JoinTime.After(c[j].JoinTime)
		}
	}

	//都没有参与打卡的
	if !c[i].JoinBool && !c[j].JoinBool {
		if c[i].Checkin.Weight != c[j].Checkin.Weight {
			return c[i].Checkin.Weight > c[j].Checkin.Weight
		}

		//权重相同的
		if c[i].Checkin.Weight == c[j].Checkin.Weight && c[i].Checkin.Weight > 0 {
			return c[i].Checkin.CreateAt.After(*(c[j].Checkin.CreateAt))
		} else if c[i].Checkin.Weight == 0 && c[j].Checkin.Weight == 0 { //没有权重
			if c[i].Checkin.JoinNum > 0 && c[j].Checkin.JoinNum > 0 {
				return c[i].Rank < c[j].Rank
			} else if c[i].Checkin.JoinNum == 0 && c[j].Checkin.JoinNum == 0 {
				return c[i].Checkin.CreateAt.After(*c[j].Checkin.CreateAt)
			} else if c[i].Checkin.JoinNum > 0 && c[j].Checkin.JoinNum == 0 {
				return c[i].Rank < c[j].Rank
			} else {
				return c[i].Rank > c[j].Rank
			}
		}
	}
	return false
}

func (c CheckinSlice) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
