package config

import (
	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile("config.yaml")
	err := viper.ReadInConfig() // 搜索并读取配置文件
	if err != nil {             // 处理错误
		log.Fatalf("没有配置文件呢？: %s \n", err)
	}
}
