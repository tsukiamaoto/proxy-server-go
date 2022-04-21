package config

import (
	"fmt"

	"github.com/spf13/viper"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	ServerAddress string
	AllowOrigins  []string
	Redis         *Redis
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
	viper.AddConfigPath("../")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Error("Config file not found")
		} else {
			// Config file was found but another error was produced
			panic("讀取設定檔出現錯誤，錯誤的原因為" + err.Error())
		}
	}

	serverAddress := fmt.Sprintf("%s:%d", viper.GetString("application.host"), viper.GetInt("application.port"))
	redis := &Redis{
		Address:  viper.GetString("redis.host"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	}
	allowOrigins := viper.GetStringSlice("application.cors.allowOrigins")

	config := &Config{
		ServerAddress: serverAddress,
		Redis:        redis,
		AllowOrigins: allowOrigins,
	}

	return config
}
