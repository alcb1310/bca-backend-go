package routes

import (
	"net/http"

	"github.com/alcb1310/bca-backend-go/internals/models"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type Router struct {
	*mux.Router

	db *gorm.DB
}

func NewRouter() *Router {
	r := &Router{
		Router: mux.NewRouter(),
	}

	r.db = models.Connect()
	r.routes()
	return r
}

func (s *Router) routes() {
	// public routes
	s.HandleFunc("/", s.handleHomeRoute()).Methods(http.MethodGet)
	s.HandleFunc("/register", s.handleRegisterRoute()).Methods(http.MethodPost)
	s.HandleFunc("/login", s.handleLogin()).Methods(http.MethodPost)

	s.initAuthRoutes()
}
