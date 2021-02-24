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
// /post:
//	 post:
//	       tags:
//	         - developers
//	       summary: insert post object
//	       requestBody:
//	         description: post object
//	         required: true
//	         content:
//	           application/json:
//	             schema:
//	               $ref: '#/components/schemas/Post'
//	       operationId: insertPost
//	       description: insert post object
//	       responses:
//	         '201':
//	           description: Information stored successfully
//	         '400':
//	           description: 'invalid input, object invalid'
//	         '500':
//	           description: service error
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

// GetPosts return posts objects
// /post:
//     get:
//       tags:
//         - developers
//       summary: return list of posts by post_name or/and author name
//       operationId: get posts
//       description: |
//         By passing in the appropriate options, get posts objects
//       parameters:
//         - in: query
//           name: post_name
//           description: post_name for searching
//           required: false
//           schema:
//             type: string
//         - in: query
//           name: author
//           description: author for searching
//           required: false
//           schema:
//             type: string
//         - in: query
//           name: order
//           description: does need ordering
//           required: false
//           schema:
//             type: boolean
//       responses:
//         '200':
//           description: search results matching criteria
//           content:
//             application/json:
//               schema:
//                 type: array
//                 items:
//                   $ref: '#/components/schemas/Post'
//         '400':
//           description: bad input parameter
//         '500':
//           description: service error
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
	}
	responce, err = json.Marshal(posts)
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

// GetPostsByAuthor return posts objects
// /post/{author}:
//     get:
//       tags:
//         - developers
//       summary: return list of posts by author name
//       operationId: getPostsByAuthor
//       description: |
//         By passing in the appropriate options, get posts objects
//       parameters:
//         - in: path
//           name: author
//           description: author for searching
//           required: true
//           schema:
//             type: string
//         - in: query
//           name: order
//           description: does need ordering
//           required: false
//           schema:
//             type: boolean
//       responses:
//         '200':
//           description: search results matching criteria
//           content:
//             application/json:
//               schema:
//                 type: array
//                 items:
//                   $ref: '#/components/schemas/Post'
//         '400':
//           description: bad input parameter
//         '500':
//           description: service error
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
	}
	responce, err = json.Marshal(posts)
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
