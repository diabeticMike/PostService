package post

import (
	"errors"
	"testing"

	"github.com/PostService/mocks"
	"github.com/PostService/model"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestInsertPost(t *testing.T) {
	t.Run("insert post with name error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cacheMock := mocks.NewMockPostCache(mockCtrl)
		post := model.Post{Name: "name1", Author: "author1"}
		postString := `{"post_name":"name1","date":"0001-01-01T00:00:00Z","author":"author1"}`
		key := "name1"
		payloadErr := errors.New("insert error")
		cacheMock.EXPECT().InsertPost(key, postString).Return(payloadErr)

		s := NewPostService(cacheMock)
		err := s.InsertPost(post)
		assert.Equal(t, payloadErr, err)
	})
	t.Run("insert post with author error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cacheMock := mocks.NewMockPostCache(mockCtrl)
		post := model.Post{Name: "name1", Author: "author1"}
		postString := `{"post_name":"name1","date":"0001-01-01T00:00:00Z","author":"author1"}`
		payloadErr := errors.New("insert error")
		cacheMock.EXPECT().InsertPost(post.Name, postString).Return(nil)
		cacheMock.EXPECT().InsertPost(post.Author, postString).Return(payloadErr)

		s := NewPostService(cacheMock)
		err := s.InsertPost(post)
		assert.Equal(t, payloadErr, err)
	})
	t.Run("success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cacheMock := mocks.NewMockPostCache(mockCtrl)
		post := model.Post{Name: "name1", Author: "author1"}
		postString := `{"post_name":"name1","date":"0001-01-01T00:00:00Z","author":"author1"}`
		cacheMock.EXPECT().InsertPost(post.Name, postString).Return(nil)
		cacheMock.EXPECT().InsertPost(post.Author, postString).Return(nil)

		s := NewPostService(cacheMock)
		err := s.InsertPost(post)
		assert.Equal(t, nil, err)
	})
}

func TestGetPostsByKey(t *testing.T) {
	t.Run("get posts by key error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cacheMock := mocks.NewMockPostCache(mockCtrl)
		key := "key"
		payloadErr := errors.New("get error")
		cacheMock.EXPECT().GetPostsByKey(key).Return(nil, payloadErr)

		s := NewPostService(cacheMock)
		posts, err := s.GetPostsByKey(key)
		assert.Equal(t, payloadErr, err)
		assert.Equal(t, posts == nil, true)
	})
	t.Run("unmarshal error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cacheMock := mocks.NewMockPostCache(mockCtrl)
		key := "key"
		cacheMock.EXPECT().GetPostsByKey(key).Return(
			[]string{`{"post_name":"name1","date":"0001`},
			nil,
		)

		s := NewPostService(cacheMock)
		posts, err := s.GetPostsByKey(key)
		assert.Equal(t, "unexpected end of JSON input", err.Error())
		assert.Equal(t, posts == nil, true)
	})
	t.Run("success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cacheMock := mocks.NewMockPostCache(mockCtrl)
		key := "name1"
		post := model.Post{Name: "name1", Author: "author1"}
		cacheMock.EXPECT().GetPostsByKey(key).Return(
			[]string{`{"post_name":"name1","date":"0001-01-01T00:00:00Z","author":"author1"}`},
			nil,
		)

		s := NewPostService(cacheMock)
		posts, err := s.GetPostsByKey(key)
		assert.Nil(t, err)
		assert.NotEmpty(t, posts)
		assert.Equal(t, posts[0], post)
	})
}

func TestGetPostsByNameAndAuthor(t *testing.T) {
	t.Run("get posts by key error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cacheMock := mocks.NewMockPostCache(mockCtrl)
		name := "name1"
		payloadErr := errors.New("get error")
		cacheMock.EXPECT().GetPostsByKey(name).Return(nil, payloadErr)

		s := NewPostService(cacheMock)
		posts, err := s.GetPostsByNameAndAuthor(name, "")
		assert.Equal(t, payloadErr, err)
		assert.Equal(t, posts == nil, true)
	})
	t.Run("unmarshal error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cacheMock := mocks.NewMockPostCache(mockCtrl)
		name := "name1"
		cacheMock.EXPECT().GetPostsByKey(name).Return(
			[]string{`{"post_name":"name1","date":"0001`},
			nil,
		)

		s := NewPostService(cacheMock)
		posts, err := s.GetPostsByNameAndAuthor(name, "")
		assert.Equal(t, "unexpected end of JSON input", err.Error())
		assert.Equal(t, posts == nil, true)
	})
	t.Run("success", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		cacheMock := mocks.NewMockPostCache(mockCtrl)
		name := "name1"
		author := "author1"
		post := model.Post{Name: "name1", Author: "author1"}
		cacheMock.EXPECT().GetPostsByKey(name).Return(
			[]string{`{"post_name":"name1","date":"0001-01-01T00:00:00Z","author":"author1"}`},
			nil,
		)

		s := NewPostService(cacheMock)
		posts, err := s.GetPostsByNameAndAuthor(name, author)
		assert.Nil(t, err)
		assert.NotEmpty(t, posts)
		assert.Equal(t, posts[0], post)
	})
}
