package services

import (
	"github.com/c0dered273/go-musthave-diploma-tpl/internal/configs"
	"github.com/c0dered273/go-musthave-diploma-tpl/internal/repositories"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

type OrdersService interface {
}

type OrdersServiceImpl struct {
	validator *validator.Validate
	cfg       *configs.ServerConfig
	usersRepo repositories.UserRepository
	logger    zerolog.Logger
}
