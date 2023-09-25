package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/alcb1310/bca-backend-go/internals/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

var (
	ErrExpiredToken = errors.New("Token has expired")
	ErrInvalidToken = errors.New("Invalid token")
)

const minSecretKeySize = 8

type Payload struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	CompanyId  uuid.UUID `json:"company_id"`
	IsLoggedIn bool      `json:"is_logged_in"`
	IssuedAt   time.Time `json:"issued_at"`
	ExpireAt   time.Time `json:"expire_at"`
}

type JWTMaker struct {
	secretKey string
}

type Maker interface {
	CreateToken(userInfo models.User, duration time.Duration) (string, error)
	VerifyToken(token string) (*Payload, error)
}

func NewPayload(u models.User, duration time.Duration) *Payload {
	payload := &Payload{
		ID:         u.ID,
		Email:      u.Email,
		CompanyId:  u.CompanyId,
		IsLoggedIn: true,
		IssuedAt:   time.Now(),
		ExpireAt:   time.Now().Add(duration),
	}
	return payload
}

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, errors.New(fmt.Sprintf("Invalid key size: must be at least %d characters", minSecretKeySize))
	}
	return &JWTMaker{secretKey}, nil
}

func (maker *JWTMaker) CreateToken(userInfo models.User, duration time.Duration) (string, error) {
	payload := NewPayload(userInfo, duration)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodES256, payload)
	return jwtToken.SignedString([]byte(maker.secretKey))
}

func (payload *Payload) Valid() error {
	if payload.IsLoggedIn && time.Now().After(payload.ExpireAt) {
		return ErrExpiredToken
	}

	return nil
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}

		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
