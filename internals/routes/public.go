package routes

import (
	"encoding/json"
	"net/http"
	"net/mail"

	"github.com/alcb1310/bca-backend-go/internals/models"
	"golang.org/x/crypto/bcrypt"
)

const salt = 15

type registerCompany struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Employees uint8  `json:"employees"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	UserName  string `json:"userName"`
}

func (s *Router) handleHomeRoute() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Hello World")
	}
}

func (s *Router) handleRegisterRoute() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reg registerCompany
		if err := json.NewDecoder(r.Body).Decode(&reg); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if reg.Name == "" || reg.ID == "" || reg.Email == "" || reg.UserName == "" || reg.Password == "" {
			http.Error(w, "Need to provide all the information required", http.StatusBadRequest)
			return
		}
		if _, err := mail.ParseAddress(reg.Email); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if reg.Employees == 0 {
			reg.Employees = 1
		}

		tx := s.db.Begin()
		c := models.Company{
			Ruc:       reg.ID,
			Name:      reg.Name,
			Employees: reg.Employees,
		}
		result := tx.Create(&c)
		if result.Error != nil {
			tx.Rollback()
			http.Error(w, result.Error.Error(), http.StatusBadRequest)
			return
		}

		pass, err := bcrypt.GenerateFromPassword([]byte(reg.Password), salt)
		if err != nil {
			tx.Rollback()
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}
		u := models.User{
			Name:     reg.UserName,
			Email:    reg.Email,
			Password: string(pass),
			Company:  c,
		}
		result = tx.Create(&u)
		if result.Error != nil {
			tx.Rollback()
			http.Error(w, result.Error.Error(), http.StatusBadRequest)
			return
		}
		tx.Commit()

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(c)
	}
}

func (s *Router) handleLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("TODO: handle Login"))
	}
}
