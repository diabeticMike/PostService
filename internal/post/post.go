package post

import (
	"io"
	"io/ioutil"

	"github.com/PostService/internal/post/cache"
	"github.com/PostService/model"
)

// Service is interface for post logic
type Service interface {
	InsertPost(post model.Post, body io.ReadCloser) error
}

// NewPostService return realization of Service interface using cache
func NewPostService(cache cache.PostCache) Service {
	return &service{cache}
}

// service is realization of the post business logic
type service struct {
	cache cache.PostCache
}

// InsertPost use cache for storing post object
func (s *service) InsertPost(post model.Post, body io.ReadCloser) error {
	nameKey := post.Name
	authorKey := post.Author
	postBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	if err := s.cache.InsertPost(nameKey, string(postBytes)); err != nil {
		return err
	}
	if err := s.cache.InsertPost(authorKey, string(postBytes)); err != nil {
		return err
	}
	return nil
}
