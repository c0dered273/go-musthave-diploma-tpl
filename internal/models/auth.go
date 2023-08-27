package models

import "github.com/golang-jwt/jwt/v4"

//go:generate easyjson auth.go
//easyjson:json
type LoginRequestDTO struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (l *LoginRequestDTO) ToUser() *User {
	return &User{
		Username: l.Login,
		Password: l.Password,
	}
}

//easyjson:json
type AuthResponseDTO struct {
	AccessToken string `json:"access_token"`
}

type AuthClaim struct {
	jwt.RegisteredClaims
}
