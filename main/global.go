// @author zhuweitung 2023/1/24
package main

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// 存储全局变量
var (
	CONFIG_FILE = "./config/config.yaml"
	CONFIG      Config
)

func LoadConfig() {
	viper.SetConfigType("yaml")
	viper.SetConfigFile(CONFIG_FILE)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&CONFIG); err != nil {
		panic(err)
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		if err := viper.Unmarshal(&CONFIG); err != nil {
			panic(err)
		}
	})
}

func SaveConfig(key string, value interface{}) {
	viper.Set(key, value)
	if err := viper.WriteConfig(); err != nil {
		panic(err)
	}
}
