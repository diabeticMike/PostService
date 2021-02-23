package post

import (
	"encoding/json"

	"github.com/PostService/internal/post/cache"
	"github.com/PostService/model"
)

// Service is interface for post logic
type Service interface {
	InsertPost(post model.Post) error
	GetPostsByKey(key string) ([]model.Post, error)
	GetPostsByNameAndAuthor(name, author string) ([]model.Post, error)
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
func (s *service) InsertPost(post model.Post) error {
	nameKey := post.Name
	authorKey := post.Author
	postBytes, err := json.Marshal(post)
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

// GetPostsByKey return posts by key(author or post name)
func (s *service) GetPostsByKey(key string) ([]model.Post, error) {
	resList, err := s.cache.GetPostsByKey(key)
	if err != nil {
		return nil, err
	}
	postList := []model.Post{}
	for _, v := range resList {
		post := model.Post{}
		if err := json.Unmarshal([]byte(v), &post); err != nil {
			return nil, err
		}
		postList = append(postList, post)
	}
	return postList, nil
}

// GetPostsByNameAndAuthor return posts by author and post name
func (s *service) GetPostsByNameAndAuthor(name, author string) ([]model.Post, error) {
	resList, err := s.cache.GetPostsByKey(name)
	if err != nil {
		return nil, err
	}

	postList := []model.Post{}
	for _, v := range resList {
		post := model.Post{}
		if err := json.Unmarshal([]byte(v), &post); err != nil {
			return nil, err
		}
		if post.Author == author {
			postList = append(postList, post)
		}
	}
	return postList, nil
}
