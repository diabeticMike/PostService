package router

import (
	"net/http"

	"github.com/PostService/infrastructure/logger"
	"github.com/PostService/internal/post"
	postCache "github.com/PostService/internal/post/cache"
	"github.com/PostService/web/controller"
	"github.com/go-redis/redis"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// New base router
func New(log logger.Logger, rc *redis.Client) (router *mux.Router,
	headers handlers.CORSOption,
	methods handlers.CORSOption,
	origins handlers.CORSOption,
	err error) {
	router = mux.NewRouter().StrictSlash(true)

	postSvc := post.NewPostService(postCache.NewPostCache(rc))
	postCntr := controller.NewPostController(log, postSvc)
	router.HandleFunc("/post", postCntr.InsertPost).Methods(http.MethodPost)
	router.HandleFunc("/post", postCntr.GetPosts).Methods(http.MethodGet)
	router.HandleFunc("/post/{author}", postCntr.GetPostsByAuthor).Methods(http.MethodGet)

	headers = handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	methods = handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins = handlers.AllowedOrigins([]string{"*"})
	return router, headers, methods, origins, nil
}
