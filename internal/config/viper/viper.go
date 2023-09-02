package viper

import (
	"github.com/spf13/viper"
)

func init() {
	// viper.SetConfigName("config")
	// viper.SetConfigType("yml")
	// Using env variables for config file path as render does not support yml files as env
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath("../../etc/secrets")

	if err := viper.ReadInConfig(); err != nil {
		panic("Error reading config file: " + err.Error())
	}
}

func Get(key string) string {
	return viper.GetString(key)
}
