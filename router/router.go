package router

import (
	"net/http"
	"nimblestack/database"
	"nimblestack/router/apis"
	"nimblestack/router/middleware"
)

type Router struct {
	queries   *database.Queries
	jwtSecret []byte
	mux       *http.ServeMux
}

func NewRouter(queries *database.Queries, jwtSecret []byte) *Router {
	r := &Router{
		queries:   queries,
		jwtSecret: jwtSecret,
		mux:       http.NewServeMux(),
	}
	r.setupRoutes()
	return r
}

func (r *Router) setupRoutes() {
	// Auth
	authHandler := apis.NewAuthApi(r.queries, r.jwtSecret)
	r.mux.HandleFunc("/api/register", authHandler.Register)
	r.mux.HandleFunc("/api/login", authHandler.Login)

	// Protected
	userHandler := apis.NewUserApi(r.queries)
	r.mux.HandleFunc("/api/me", middleware.CheckAuth(userHandler.GetCurrentUser, r.jwtSecret))
}

func (r *Router) Handler() http.Handler {
	return r.mux
}
