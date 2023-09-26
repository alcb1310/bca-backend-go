package routes

import (
	"net/http"
)

func (p *protectedRoutes) createUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("create user"))
	}
}

func (p *protectedRoutes) getAllUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("get all users"))
	}
}

func (p *protectedRoutes) getOneUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("get one user"))
	}
}

func (p *protectedRoutes) updateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("update user"))
	}
}

func (p *protectedRoutes) deleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("delete user"))
	}
}
