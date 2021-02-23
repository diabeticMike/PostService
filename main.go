package main

import (
	baseLog "log"
	"net/http"
	"strings"

	"github.com/PostService/infrastructure/config"
	"github.com/PostService/infrastructure/datastore"
	"github.com/PostService/infrastructure/logger"
	"github.com/PostService/web/router"
	"github.com/go-redis/redis"
	"github.com/gorilla/handlers"
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

	// Create service logger
	if log, err = logger.New(conf.Log); err != nil {
		baseLog.Fatal(err.Error())
	}

	if redisClient, err = datastore.NewRedis(redis.Options{
		Addr:     conf.Redis.Address,
		Password: conf.Redis.Password,
		DB:       conf.Redis.DB,
	}); err != nil {
		log.Fatal(err.Error())
	}

	requestInfo := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := strings.Split(r.RemoteAddr, ":")[0]
			message := " | " + ip + " | " + r.Method + " | " + r.URL.RequestURI()
			log.Info("main", message)
			h.ServeHTTP(w, r)
		})
	}

	mainRouter, headers, methods, origins, err := router.New(log, redisClient)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Fatal(http.ListenAndServe(conf.ListenPort, requestInfo(handlers.CORS(headers, methods, origins)(mainRouter))))
}
