package model

// APIResponse 统一响应结构
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type ResponseCheckinItem struct {
	Id            int    `json:"id"`
	Title         string `json:"title"`
	CreateAt      string `json:"createAt"`
	UpdateAt      string `json:"updateAt"`
	CheckinStatus int    `json:"checkinStatus"`
	JoinBool      bool   `json:"joinBool"`
	RecordBool    bool   `json:"recordBool"`
	JoinNumber    int64  `json:"joinNumber"`
	Rank          int    `json:"rank"`
	Weight        int    `json:weight`
	JoinTime      string `json:"joinTime"`
}

// type ResponseCheckinJoin struct {
// 	Id       int    `json:"id"`
// 	Uid      int    `json:"uid"`
// 	Cid      int    `json:"cid"`
// 	JoinTime string `json:"joinTime"`
// 	CreateAt string `json:"createAt"`
// 	UpdateAt string `json:"updateAt"`
// 	Status   int    `json:"status"`
// 	JoinBool bool   `json:joinBool`
// }
// type ResponseCheckinRecord struct {
// 	Id         int    `json:"id"`
// 	Uid        int    `json:"uid"`
// 	Cid        int    `json:"cid"`
// 	Date       int    `json:"date"`
// 	CreateAt   string `json:"createAt"`
// 	UpdateAt   string `json:"updateAt"`
// 	Status     int    `json:"status"`
// 	RecordBool bool   `json:recordBool`
// }
