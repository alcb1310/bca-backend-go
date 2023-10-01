package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/mail"

	"github.com/alcb1310/bca-backend-go/internals/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type userWithoutPassword struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	CompanyId uuid.UUID `json:"company_id"`
}

var role models.Role

func (p *protectedRoutes) createUser() http.HandlerFunc {
	type createUserJSON struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Role     string `json:"role"`
	}
	var (
		user    createUserJSON
		count   int64
		company models.Company
	)

	return func(w http.ResponseWriter, r *http.Request) {
		token, err := GetMyPaload(r)
		if err != nil {
			log.Println(":Error: ", err)
			return
		}

		// Data validation
		if token.Role != "admin" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Missing data", http.StatusBadRequest)
			return
		}
		if user.Name == "" || user.Email == "" || user.Password == "" || user.Role == "" {
			http.Error(w, "Need to provide all the information required", http.StatusBadRequest)
			return
		}
		if _, err := mail.ParseAddress(user.Email); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), salt)
		if err != nil {
			// tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Check for maximum number of employees
		p.db.Model(&models.User{}).Where("company_id = ?", token.CompanyId).Count(&count)
		p.db.Find(&company, "id = ?", token.CompanyId)
		if int64(company.Employees) < count+1 {
			http.Error(w, "Employee count exceeded", http.StatusForbidden)
			return
		}

		// Get the role id
		result := p.db.Find(&role, "name=?", user.Role)
		if result.Error != nil || result.RowsAffected != 1 {
			if res := p.db.Find(&role, "name='user'"); res.Error != nil || res.RowsAffected != 1 {
				http.Error(w, res.Error.Error(), http.StatusInternalServerError)
				return
			}
		}

		// Create the user
		u := models.User{
			Email:     user.Email,
			Name:      user.Name,
			Password:  string(pass),
			Role:      role,
			CompanyId: token.CompanyId,
		}
		result = p.db.Create(&u)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(u)
	}
}

func (p *protectedRoutes) getAllUsers() http.HandlerFunc {

	var user []userWithoutPassword

	return func(w http.ResponseWriter, r *http.Request) {
		token, err := GetMyPaload(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if token.Role == "admin" {
			p.db.Find(&user, "company_id = ?", token.CompanyId)
		} else {
			p.db.Find(&user, "company_id = ? and id = ?", token.CompanyId, token.ID)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	}
}

func (p *protectedRoutes) getOneUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := GetMyPaload(r)
		if err != nil {
			log.Println(":Error: ", err)
			return
		}
		params := mux.Vars(r)
		parsedId, err := uuid.Parse(params["userId"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		if token.Role != "admin" && parsedId != token.ID {
			http.Error(w, "You have to be an admin to view that user", http.StatusForbidden)
			return
		}

		var u userWithoutPassword
		result := p.db.Find(&u, "company_id = ? and id = ?", token.CompanyId, parsedId)
		if result.RowsAffected != 1 {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(u)
	}
}

func (p *protectedRoutes) updateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := GetMyPaload(r)
		if err != nil {
			log.Println(":Error: ", err)
			return
		}

		fmt.Println(token)
		w.Write([]byte("update user"))
	}
}

func (p *protectedRoutes) deleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := GetMyPaload(r)
		if err != nil {
			log.Println(":Error: ", err)
			return
		}

		if token.Role != "admin" {
			http.Error(w, "You have to be an admin to delete users", http.StatusForbidden)
			return
		}
		params := mux.Vars(r)
		validId, err := uuid.Parse(params["userId"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if validId == token.ID {
			http.Error(w, "You can't delete yourself", http.StatusBadRequest)
			return
		}

		fmt.Println(token.ID, token.CompanyId)
		result := p.db.Delete(&models.User{}, "id = ? and company_id = ?", validId, token.CompanyId)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
