package service

import "demo1/internal/app/mydemo/dao/kafka"

// 连接kafka
func ServiceInitKafka() (err error) {
	return kafka.DaoInitKafka()
}

// 关闭
func Closekafka() {
	kafka.Closekafka()
}
