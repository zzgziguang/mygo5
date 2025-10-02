package model

// User 定义用户结构体
type User struct {
	ID       int    `json:"id" gorm:"primaryKey" mapstructure:"id"`
	Username string `json:"username" gorm:"not null;size:100" mapstructure:"username"`
	Email    string `json:"email" gorm:"unique;not null;size:100" mapstructure:"email"`
	Age      int    `json:"age" gorm:"check:age >= 0 AND age <= 150" mapstructure:"age"`
	Phone    int    `json:"phone" gorm:"type:int" mapstructure:"phone"`
	//Tianqi   string `json:"tianqi" gorm:"type:string" mapstructure:"tianqi"`
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
