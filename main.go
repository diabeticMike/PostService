package main

import (
	"fmt"
	baseLog "log"

	"github.com/PostService/infrastructure/config"
	"github.com/PostService/infrastructure/datastore"
	"github.com/PostService/infrastructure/logger"
	"github.com/PostService/internal/post/cache"
	"github.com/go-redis/redis"
)

func main() {
	configFilePath := "config.json"
	var (
		conf        config.Configuration
		log         logger.Logger
		redisClient *redis.Client
		err         error
	)

	// Create service configuration
	if conf, err = config.New(configFilePath); err != nil {
		baseLog.Fatal(err.Error())
	}
	fmt.Println(conf)

	// Create service logger
	if log, err = logger.New(conf.Log); err != nil {
		baseLog.Fatal(err.Error())
	}
	fmt.Println(log)

	if redisClient, err = datastore.NewRedis(redis.Options{
		Addr:     conf.Redis.Address,
		Password: conf.Redis.Password,
		DB:       conf.Redis.DB,
	}); err != nil {
		log.Fatal(err.Error())
	}
	pIntf := cache.NewPostCache(redisClient)
	if err != nil {
		log.Fatal(err.Error())
	}
	resp, err := pIntf.GetPostsByKey("shit")
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(resp)
}
