package config

import (
	"github.com/spf13/viper"
	"log"
)

var config *viper.Viper

func init() {
	var err error
	config = viper.New()
	config.SetConfigType("yaml")
	config.SetConfigName("config")
	config.AddConfigPath("./config")
	config.AddConfigPath("../config/")
	err = config.ReadInConfig()
	if err != nil {
		log.Fatal("read config error: ", err)
	}
}

func GetConfig() *viper.Viper {
	return config
}
