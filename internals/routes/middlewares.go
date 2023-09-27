package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/alcb1310/bca-backend-go/internals/models"
	"github.com/alcb1310/bca-backend-go/internals/utils"
	"github.com/google/uuid"
)

type errorResponse struct {
	Error string `json:"error"`
}

func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (s *Router) authVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearerToken := r.Header.Get("Authorization")
		if bearerToken == "" {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(&errorResponse{
				Error: "Missing autherization token",
			})
			return
		}

		token := strings.Split(bearerToken, " ")
		if len(token) != 2 {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(&errorResponse{
				Error: "Invalid autherization token",
			})
			return
		}

		secretKey := os.Getenv("SECRET")
		maker, err := utils.NewJWTMaker(secretKey)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(&errorResponse{
				Error: "Unable to authenticate",
			})
			return
		}

		tokenData, err := maker.VerifyToken(token[1])
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(&errorResponse{
				Error: "Unable to authenticate",
			})
			return
		}

		marshalStr, _ := json.Marshal(tokenData)
		ctx := r.Context()
		ctx = context.WithValue(r.Context(), "token", marshalStr)
		r = r.Clone(ctx)
		var u models.LoggedInUser
		result := s.db.First(&u, "email=?", tokenData.Email)
		if result.Error != nil || result.RowsAffected != 1 {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(&errorResponse{
				Error: "Invalid token",
			})
			return
		}

		if !bytes.Equal([]byte(token[1]), u.JWT) {
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(&errorResponse{
				Error: "Invalid token",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}

type contextPayload struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	CompanyId  uuid.UUID `json:"company_id"`
	Role       string    `json:"role"`
	IsLoggedIn bool      `json:"is_logged_in"`
	IssuedAt   time.Time `json:"issued_at"`
	ExpiredAt  time.Time `json:"expired_at"`
}

func SetMyPayload(r *http.Request, p contextPayload) {}

func GetMyPaload(r *http.Request) (contextPayload, error) {
	ctx := r.Context()
	val := ctx.Value("token")

	x, ok := val.([]byte)
	if !ok {
		return contextPayload{}, errors.New("Unable to load context")
	}
	var p contextPayload
	err := json.Unmarshal(x, &p)
	if err != nil {
		return contextPayload{}, errors.New("Unable to parse context")
	}
	return p, nil
}
