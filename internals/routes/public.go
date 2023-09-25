package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"os"
	"time"

	"github.com/alcb1310/bca-backend-go/internals/models"
	"github.com/alcb1310/bca-backend-go/internals/utils"
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

type loginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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
		var c loginCredentials
		var u models.User

		if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		result := s.db.Find(&u, "email = ?", c.Email)
		if result.Error != nil || result.RowsAffected != 1 {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(c.Password)); err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		secretKey := os.Getenv("SECRET")
		jwtMaker, err := utils.NewJWTMaker(secretKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		token, err := jwtMaker.CreateToken(u, 60*time.Minute)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		byteToken := []byte(token)
		loggedInUser := models.LoggedInUser{
			Email: u.Email,
			JWT:   byteToken,
		}
		s.db.Save(&loggedInUser)

		w.WriteHeader(http.StatusOK)
		type response struct {
			Response string `json:"response"`
		}
		json.NewEncoder(w).Encode(response{
			Response: fmt.Sprintf("Bearer %s", token),
		})
	}
}
