package services

import (
	"context"

	"github.com/c0dered273/go-musthave-diploma-tpl/internal/models"
	"github.com/c0dered273/go-musthave-diploma-tpl/internal/store"
	"github.com/rs/zerolog"
)

type HealthService interface {
	ConnPing(ctx context.Context) error
}

type HealthServiceImpl struct {
	connCheck store.ConnCheck
	logger    zerolog.Logger
}

func NewHealthService(logger zerolog.Logger, connCheck store.ConnCheck) HealthService {
	return &HealthServiceImpl{
		connCheck: connCheck,
		logger:    logger,
	}
}

func (h HealthServiceImpl) ConnPing(ctx context.Context) error {
	err := h.connCheck.Ping(ctx)
	if err != nil {
		return models.NewErrInternal(err, "DB_ERROR", "Connection check failed")
	}
	return nil
}
