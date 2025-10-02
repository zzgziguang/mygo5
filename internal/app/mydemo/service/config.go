package service

import (
	"demo1/internal/app/mydemo/conf"
	"os"

	"gopkg.in/yaml.v2"
)

var Cfg *conf.Config

// 加载配置文件
func LoadConfig() (err error) {
	Cfg = &conf.Config{}
	data, err := os.ReadFile("../../../configs/config.yaml")
	if err != nil {
		return
	}

	err = yaml.Unmarshal(data, Cfg)

	return
}
