package main

import (
	"encoding/json"
	Config "tsukiamaoto/proxy-server-go/config"
	"tsukiamaoto/proxy-server-go/proxy"
	"tsukiamaoto/proxy-server-go/redis"

	"github.com/gin-gonic/gin"
	"github.com/jasonlvhit/gocron"
	log "github.com/sirupsen/logrus"
)

var redisDB *redis.Redis

type DataResponse struct {
	Data []string `json:"data"`
}

func init() {
	redisDB = redis.New()
}

func main() {
	// fetch proxy list
	go autoFetchProxy()
	// load config
	config := Config.LoadConfig()

	server := gin.Default()
	server.GET("/api/v1/proxy", func(c *gin.Context) {
		var proxies []string
		res := redisDB.JSONGet("proxy", ".")
		err := json.Unmarshal(res.([]byte), &proxies)
		if err != nil {
			log.Error("Failed to JSON Unmarshal proxies: ", err)
		}

		c.JSON(200, DataResponse{Data: proxies})

	})

	server.Run(config.ServerAddress)
}

func autoFetchProxy() {
	schedule := gocron.NewScheduler()
	schedule.Every(1).Minute().Do(proxy.FetchTask)
	<-schedule.Start()
}
