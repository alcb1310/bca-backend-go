package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Router struct {
	*mux.Router
}

func NewRouter() *Router {
	r := &Router{
		Router: mux.NewRouter(),
	}
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
