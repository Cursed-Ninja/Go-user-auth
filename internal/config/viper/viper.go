package viper

import (
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../../etc/secrets")

	if err := viper.ReadInConfig(); err != nil {
		panic("Error reading config file: " + err.Error())
	}
}

func Get(key string) string {
	return viper.GetString(key)
}
