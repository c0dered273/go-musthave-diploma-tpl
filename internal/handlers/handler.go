package handlers

import (
	"errors"
	"net/http"

	"github.com/c0dered273/go-musthave-diploma-tpl/internal/entities"
	"github.com/c0dered273/go-musthave-diploma-tpl/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/rs/zerolog"
)

func NewHandler(logger zerolog.Logger, services services.ServiceContext) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	//r.Use(middleware.RequestID)
	r.Use(httplog.RequestLogger(logger))
	//r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.NotFound(notFound(logger))
	r.MethodNotAllowed(notAllowed(logger))

	r.Route("/health", func(r chi.Router) {
		r.Get("/livez", liveProbe(logger))
		r.Get("/readyz", readyProbe(logger, services.HealthService))
	})

	return r
}

func notAllowed(logger zerolog.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := entities.NewErrNotAllowed(
			errors.New("method not allowed"),
			"HTTP_ERROR",
			"Method not allowed",
		)
		wsErr := entities.WriteStatusError(w, err)
		if wsErr != nil {
			logger.Error().Err(err)
			return
		}
	}
}

func notFound(logger zerolog.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := entities.NewErrNotFound(
			errors.New("endpoint not found"),
			"HTTP_ERROR",
			"Endpoint not found",
		)
		wsErr := entities.WriteStatusError(w, err)
		if wsErr != nil {
			logger.Error().Err(err)
			return
		}
	}
}

func liveProbe(logger zerolog.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("."))
		if err != nil {
			logger.Error().Err(err)
		}
	}
}

func readyProbe(logger zerolog.Logger, health services.HealthService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		if err := health.DBConnPing(r.Context()); err != nil {
			errSts := entities.WriteStatusError(w, err)
			if errSts != nil {
				logger.Error().Err(errSts)
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("."))
		if err != nil {
			logger.Error().Err(err)
		}
	}
}
