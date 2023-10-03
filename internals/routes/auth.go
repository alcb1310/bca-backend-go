package routes

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/alcb1310/bca-backend-go/internals/models"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type protectedRoutes struct {
	*mux.Router

	db *gorm.DB
}

func (s *Router) initAuthRoutes() {
	p := &protectedRoutes{
		db:     s.db,
		Router: s.PathPrefix("/api/v1").Subrouter(),
	}
	p.Use(jsonMiddleware)
	p.Use(s.authVerify)

	p.HandleFunc("/logout", p.handleLogout()).Methods(http.MethodGet)

	// users endpoints
	p.HandleFunc("/users", p.createUser()).Methods(http.MethodPost)
	p.HandleFunc("/users", p.getAllUsers()).Methods(http.MethodGet)
	p.HandleFunc("/users/{userId}", p.getOneUser()).Methods(http.MethodGet)
	p.HandleFunc("/users/{userId}", p.updateUser()).Methods(http.MethodPut)
	p.HandleFunc("/users/{userId}", p.deleteUser()).Methods(http.MethodDelete)
}

func (p *protectedRoutes) handleLogout() http.HandlerFunc {
	type response struct {
		Response string `json:"response"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		token, err := GetMyPaload(r)
		if err != nil {
			log.Println(":Error: ", err)
			return
		}

		p.db.Delete(&models.LoggedInUser{}, "email = ?", token.Email)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response{
			Response: "User logged out",
		})
	}
}
