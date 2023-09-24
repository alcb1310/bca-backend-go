package routes

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type protectedRoutes struct {
	*mux.Router

	db *gorm.DB
}

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(":INFO: Should only display on logout")
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (s *Router) initAuthRoutes() {
	p := &protectedRoutes{
		db:     s.db,
		Router: s.PathPrefix("/api/v1").Subrouter(),
	}
	p.Use(jsonMiddleware)

	p.HandleFunc("/logout", p.handleLogout()).Methods(http.MethodGet)
}

func (p *protectedRoutes) handleLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("TODO: handle Logout"))
	}
}
