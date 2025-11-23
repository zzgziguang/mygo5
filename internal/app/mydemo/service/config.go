package service

import (
	"demo1/internal/app/mydemo/conf"
	"flag"
	"os"

	"gopkg.in/yaml.v2"
)

var Cfg *conf.Config

// 加载配置文件
func LoadConfig() (err error) {
	configFile := flag.String("conf", "", "config文件")

	flag.Parse()

	Cfg = &conf.Config{}
	data, err := os.ReadFile(*configFile)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(data, Cfg)
	return
}
