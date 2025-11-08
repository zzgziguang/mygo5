package model

// APIResponse 统一响应结构
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type WeatherItem struct {
	Temperature string `json:"temperature"`
	Weather     string `json:"weather"`
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
	Weight        int    `json:"weight"`
	JoinTime      string `json:"joinTime"`
	CidNum        int    `json:"cidNum"`
}
