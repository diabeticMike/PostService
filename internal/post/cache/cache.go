package cache

import (
	"github.com/go-redis/redis"
)

// PostCache used for redis logic related to post entity
type PostCache interface {
	GetPostsByKey(key string) ([]string, error)
	InsertPost(key, post string) error
}

// NewPostCache return new PostCache realization
func NewPostCache(rc *redis.Client) PostCache {
	return &postCache{rc}
}

type postCache struct {
	rc *redis.Client
}

func (pr *postCache) GetPostsByKey(key string) ([]string, error) {
	resp, err := pr.rc.LRange(key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (pr *postCache) InsertPost(key, post string) error {
	err := pr.rc.LPush(key, post).Err()
	if err != nil {
		return err
	}

	return nil
}
