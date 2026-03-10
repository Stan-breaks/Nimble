package router

import (
	"net/http"
	"nimblestack/database"
	"nimblestack/router/apis"
	"nimblestack/router/middleware"
)

type Router struct {
	queries   *database.Queries
	jwtSercet []byte
	mux       *http.ServeMux
}

func NewRouter(queries *database.Queries, jwtSercet []byte) *Router {
	r := &Router{
		queries:   queries,
		jwtSercet: jwtSercet,
		mux:       http.NewServeMux(),
	}
	r.setupRoutes()
	return r
}

func (r *Router) setupRoutes() {
	// serving static files
	fs := http.FileServer(http.Dir("public"))
	r.mux.Handle("/static/", http.StripPrefix("/static/", middleware.AddContentType(fs)))

// apis section
	authHander := apis.NewAuthApi(r.queries, r.jwtSercet)
	r.mux.HandleFunc("/api/login", authHander.Login)
	r.mux.HandleFunc("/api/register", authHander.Reqister)
	// serving protected apis
	userHandler := apis.NewUserApi(r.queries)
	r.mux.HandleFunc("/api/me", middleware.CheckAuth(userHandler.GetCurrentUSer, r.jwtSercet))
	dashHandler := apis.NewDashApi(r.queries)
	r.mux.HandleFunc("/api/coordinator/students", middleware.CheckAuth(dashHandler.GetAllStudents, r.jwtSercet))
	r.mux.HandleFunc("/api/coordinator/projects", middleware.CheckAuth(dashHandler.GetProjectsData, r.jwtSercet))
	r.mux.HandleFunc("/api/coordinator/supervisors", middleware.CheckAuth(dashHandler.GetAllSupervisors, r.jwtSercet))
	r.mux.HandleFunc("/api/coordinator/assign", middleware.CheckAuth(dashHandler.AssignSupervisor, r.jwtSercet))
}

func (r *Router) Handler() http.Handler {
	return r.mux
}
