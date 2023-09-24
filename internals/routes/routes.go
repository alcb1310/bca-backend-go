package routes

import (
	"encoding/json"
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
	s.HandleFunc("/", s.handleHomeRoute()).Methods(http.MethodGet)
}

func (s *Router) handleHomeRoute() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Hello World")
	}
}
