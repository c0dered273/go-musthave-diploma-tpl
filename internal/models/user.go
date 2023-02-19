package models

//go:generate easyjson user.go

type User struct {
	Username string
	Password string
}

//easyjson:json
type UserBalanceDTO struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn,omitempty"`
}
