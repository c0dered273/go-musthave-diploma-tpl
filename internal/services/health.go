package services

import (
	"context"

	"github.com/c0dered273/go-musthave-diploma-tpl/internal/entities"
	"github.com/c0dered273/go-musthave-diploma-tpl/internal/repositories"
	"github.com/rs/zerolog"
)

type HealthService interface {
	DBConnPing(ctx context.Context) error
}

type HealthServiceImpl struct {
	repo   repositories.Repository
	logger zerolog.Logger
}

func NewHealthService(repo repositories.Repository, logger zerolog.Logger) HealthService {
	return HealthServiceImpl{
		repo:   repo,
		logger: logger,
	}
}

func (h HealthServiceImpl) DBConnPing(ctx context.Context) error {
	err := h.repo.Ping(ctx)
	if err != nil {
		return entities.NewErrInternal(err, "DB_ERROR", "Connection check failed")
	}
	return nil
}
