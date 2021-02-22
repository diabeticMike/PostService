package main

import (
	"fmt"
	baseLog "log"

	"github.com/PostService/infrastructure/config"
	"github.com/PostService/infrastructure/logger"
)

func main() {
	configFilePath := "config.json"
	var (
		Config config.Configuration
		log    logger.Logger
		// redisClient *redis.Client
		err error
	)

	// Create service configuration
	if Config, err = config.New(configFilePath); err != nil {
		baseLog.Fatal(err.Error())
	}
	fmt.Println(Config)

	// Create service logger
	if log, err = logger.New(Config.Log); err != nil {
		baseLog.Fatal(err.Error())
	}
}
