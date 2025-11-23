package model

type CheckInMsg struct {
	Uid       int    `json:"uid"`
	Cid       int    `json:"cid"`
	Timestamp int64  `json:"timestamp"` // Unix 时间戳
	Msg       string `json:"msg"`
}
