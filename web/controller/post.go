package controller

import (
	"encoding/json"
	"net/http"

	"github.com/PostService/infrastructure/logger"
	"github.com/PostService/internal/post"
	"github.com/PostService/model"
)

// NewPostController return PostController instance by passing log and post's business logic interface
func NewPostController(log logger.Logger, postSvc post.Service) *PostController {
	return &PostController{log: log, postSvc: postSvc}
}

// PostController responsible for holding logger and interface for post business logic
type PostController struct {
	log     logger.Logger
	postSvc post.Service
}

// InsertPost create post record
func (pc *PostController) InsertPost(w http.ResponseWriter, r *http.Request) {
	var post model.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		pc.log.Error(err.Error())
		return
	}
	if err := pc.postSvc.InsertPost(post); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		pc.log.Error(err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(`Information stored successfully`))
	w.WriteHeader(http.StatusOK)
}

// GetPosts posts objects
func (pc *PostController) GetPosts(w http.ResponseWriter, r *http.Request) {
	qParams := r.URL.Query()
	author := qParams.Get("author")
	name := qParams.Get("post_name")
	var (
		err   error
		posts []model.Post
	)
	if author != "" && name != "" {
		posts, err = pc.postSvc.GetPostsByNameAndAuthor(name, author)
		if err != nil {
			pc.log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else if author != "" || name != "" {
		var key string
		if author != "" {
			key = author
		}
		if name != "" {
			key = name
		}
		if posts, err = pc.postSvc.GetPostsByKey(key); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			pc.log.Error(err.Error())
			return
		}
	}

	responce, err := json.Marshal(posts)
	if err != nil {
		pc.log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(responce); err != nil {
		pc.log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
