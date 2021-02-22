package controller

import (
	"encoding/json"
	"net/http"

	"github.com/PostService/infrastructure/logger"
	"github.com/PostService/internal/post"
	"github.com/PostService/model"
)

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
	if err := pc.postSvc.InsertPost(post, r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		pc.log.Error(err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(`Information stored successfully`))
	w.WriteHeader(http.StatusOK)
}
