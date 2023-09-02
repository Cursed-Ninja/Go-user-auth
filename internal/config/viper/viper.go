package viper

import (
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("../../")
	viper.AddConfigPath("../../internal/config")

	if err := viper.ReadInConfig(); err != nil {
		panic("Error reading config file: " + err.Error())
	}
}

func Get(key string) string {
	return viper.GetString(key)
}
