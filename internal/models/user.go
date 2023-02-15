package models

//go:generate easyjson user.go

type User struct {
	Username string
	Password string
}

type Users []User
