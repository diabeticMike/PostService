package router

import (
	"database/sql"

	"github.com/PostService/infrastructure/logger"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// New base router
func New(log logger.Logger, db *sql.DB) (router *mux.Router,
	headers handlers.CORSOption,
	methods handlers.CORSOption,
	origins handlers.CORSOption,
	err error) {
	router = mux.NewRouter().StrictSlash(true)

	headers = handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})
	methods = handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins = handlers.AllowedOrigins([]string{"*"})
	return router, headers, methods, origins, nil
}
