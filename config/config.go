package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	AllowOrigins []string
	Redis        *Redis
}

type Redis struct {
	Address  string
	Password string
	DB       int
}

type Database struct {
	Name   string
	Source string
}

func LoadConfig() *Config {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		panic("讀取設定檔出現錯誤，錯誤的原因為" + err.Error())
	}

	redis := &Redis{
		Address:  viper.GetString("redis.host"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	}
	allowOrigins := viper.GetStringSlice("application.cors.allowOrigins")

	config := &Config{
		Redis:        redis,
		AllowOrigins: allowOrigins,
	}

	return config
}