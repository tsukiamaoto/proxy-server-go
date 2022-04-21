package main

import (
	Config "tsukiamaoto/proxy-server-go/config"
	"tsukiamaoto/proxy-server-go/proxy"
	"tsukiamaoto/proxy-server-go/redis"

	"github.com/gin-gonic/gin"
	"github.com/jasonlvhit/gocron"
)

var redisDB *redis.Redis

type DataResponse struct {
	Data interface{} `json:"data"`
}

func init() {
	redisDB = redis.New()
	redisDB.ConnectRDB()
}

func main() {
	// fetch proxy list
	go autoFetchProxy()
	// load config
	config := Config.LoadConfig()

	server := gin.Default()
	server.GET("/api/v1/proxy", func(c *gin.Context) {
		proxies := redisDB.Get("proxy")
		c.JSON(200, DataResponse{Data: proxies})
	})

	server.Run(config.ServerAddress)
}

func autoFetchProxy() {
	schedule := gocron.NewScheduler()
	schedule.Every(1).Minute().Do(proxy.FetchTask)
	<-schedule.Start()
}
