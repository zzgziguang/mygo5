package model

// User 定义用户结构体
type User struct {
	Id       int    `json:"id" gorm:"column:id" mapstructure:"id"`
	Username string `json:"username" gorm:"column:username" mapstructure:"username"`
	Email    string `json:"email" gorm:"column:email" mapstructure:"email"`
	Age      int    `json:"age" gorm:"column:age" mapstructure:"age"`
	Phone    int    `json:"phone" gorm:"column:phone" mapstructure:"phone"`
	//Tianqi   string `json:"tianqi" mapstructure:"tianqi"`
}

// 指定User对应的表名
func (User) TableName() string {
	return "users"
}

// 按年龄升序排序
type ByAgeAsc []User

func (a ByAgeAsc) Len() int           { return len(a) }
func (a ByAgeAsc) Less(i, j int) bool { return a[i].Age < a[j].Age }
func (a ByAgeAsc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// 按年龄降序排序
type ByAgeDesc []User

func (a ByAgeDesc) Len() int           { return len(a) }
func (a ByAgeDesc) Less(i, j int) bool { return a[i].Age > a[j].Age }
func (a ByAgeDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
