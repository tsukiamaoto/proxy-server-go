package redis

import (
	Config "tsukiamaoto/proxy-server-go/config"

	"context"

	"github.com/go-redis/redis/v8"
	"github.com/nitishm/go-rejson/v4"
	log "github.com/sirupsen/logrus"
)

type Redis struct {
	Handler *rejson.Handler
	Client  *redis.Client
}

var ctx = context.Background()

func New() *Redis {
	config := Config.LoadConfig()
	client := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Address,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})

	handler := rejson.NewReJSONHandler()
	handler.SetGoRedisClient(client)

	return &Redis{
		Handler: handler,
		Client:  client,
	}
}

func (r *Redis) JSONGet(key, path string) interface{} {
	value, err := r.Handler.JSONGet(key, path)
	if err != nil {
		log.Error("Failed to JSONGet: ", err)
		return nil
	}
	return value
}

func (r *Redis) JSONSet(key, path string, value interface{}) {
	_, err := r.Handler.JSONSet(key, path, value)
	if err != nil {
		log.Error("Failed to JSONSet: ", err)
	}
}

func (r *Redis) Exists(key string) bool {
	ok, err := r.Client.Exists(ctx, key).Result()
	if err != nil {
		log.Error("Failed to Exists: ", err)
	}

	return ok == 1
}
