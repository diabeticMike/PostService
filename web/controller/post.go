package controller

import (
	"encoding/json"
	"net/http"
	"sort"
	"time"

	"github.com/PostService/infrastructure/logger"
	"github.com/PostService/internal/post"
	"github.com/PostService/model"
	"github.com/gorilla/mux"
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
	var post = struct {
		Name   string `json:"post_name"`
		Date   string `json:"date"`
		Author string `json:"author"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		pc.log.Error(err.Error())
		return
	}
	t, err := time.Parse("02.01.06", post.Date)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		pc.log.Error(err.Error())
		return
	}
	if err := pc.postSvc.InsertPost(model.Post{Name: post.Name, Date: t, Author: post.Author}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		pc.log.Error(err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	if _, err := w.Write([]byte(`Information stored successfully`)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		pc.log.Error(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

// GetPosts posts objects
func (pc *PostController) GetPosts(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		posts model.Posts
	)
	qParams := r.URL.Query()
	author := qParams.Get("author")
	name := qParams.Get("post_name")
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
		} else if name != "" {
			key = name
		}
		if posts, err = pc.postSvc.GetPostsByKey(key); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			pc.log.Error(err.Error())
			return
		}
	}
	if len(posts) == 0 {
		w.Header().Set("Content-Type", "text/plain")
		if _, err := w.Write([]byte(`No posts found`)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			pc.log.Error(err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	order := qParams.Get("order")
	responce := []byte{}
	if order == "true" {
		sort.Sort(&posts)
		responce, err = json.Marshal(posts)
		if err != nil {
			pc.log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		responce, err = json.Marshal(posts)
		if err != nil {
			pc.log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(responce); err != nil {
		pc.log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetPosts posts objects
func (pc *PostController) GetPostsByAuthor(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		posts model.Posts
	)
	author, ok := mux.Vars(r)["author"]
	if !ok {
		w.Header().Set("Content-Type", "text/plain")
		if _, err := w.Write([]byte(`No author provided`)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			pc.log.Error(err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	posts, err = pc.postSvc.GetPostsByKey(author)
	if err != nil {
		pc.log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(posts) == 0 {
		w.Header().Set("Content-Type", "text/plain")
		if _, err := w.Write([]byte(`No posts found`)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			pc.log.Error(err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}

	qParams := r.URL.Query()
	order := qParams.Get("order")
	responce := []byte{}
	if order == "true" {
		sort.Sort(&posts)
		responce, err = json.Marshal(posts)
		if err != nil {
			pc.log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		responce, err = json.Marshal(posts)
		if err != nil {
			pc.log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(responce); err != nil {
		pc.log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
