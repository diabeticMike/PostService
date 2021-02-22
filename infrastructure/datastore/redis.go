package datastore

import (
	"github.com/go-redis/redis"
)

// NewRedis connects to redis and return redis client
func NewRedis(opts redis.Options) (client *redis.Client, err error) {
	client = redis.NewClient(&opts)
	if _, err = client.Ping().Result(); err != nil {
		return nil, err
	}
	return
}
