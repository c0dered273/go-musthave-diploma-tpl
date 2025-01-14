package handlers

import (
	"errors"
	"net/http"

	"github.com/c0dered273/go-musthave-diploma-tpl/internal/configs"
	middleware2 "github.com/c0dered273/go-musthave-diploma-tpl/internal/middleware"
	"github.com/c0dered273/go-musthave-diploma-tpl/internal/models"
	"github.com/c0dered273/go-musthave-diploma-tpl/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/rs/zerolog"
)

func NewHandler(logger zerolog.Logger, cfg *configs.ServerConfig, services *services.ServiceContext) http.Handler {
	httpLogger := logger.With().Str("module", "handler").Logger()
	r := chi.NewRouter()
	r.Use(middleware.AllowContentType("text/plain", "application/json"))
	r.Use(middleware.RealIP)
	//r.Use(middleware.RequestID)
	r.Use(httplog.RequestLogger(httpLogger))
	//r.Use(middleware.Recoverer)
	r.Use(middleware2.SetContentType("application/json"))
	r.Use(middleware.Compress(5))

	r.NotFound(notFound())
	r.MethodNotAllowed(notAllowed())

	if cfg.Server.PprofEnable {
		r.Mount("/debug", middleware.Profiler())
	}

	r.Route("/health", func(r chi.Router) {
		r.Get("/livez", liveProbe(httpLogger))
		r.Get("/readyz", readyProbe(httpLogger, services.HealthService))
	})

	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", registerUser(httpLogger, services.UsersService))
		r.Post("/login", loginUser(httpLogger, services.UsersService))

		r.Group(func(r chi.Router) {
			r.Use(middleware2.JwtVerifier(httpLogger, cfg.APISecret))
			r.Post("/orders", addOrders(logger, services.UsersService))
			r.Get("/orders", getUserOrders(logger, services.UsersService))
			r.Get("/balance", getUserBalance(logger, services.UsersService))
			r.Post("/balance/withdraw", withdrawBalance(logger, services.UsersService))
			r.Get("/withdrawals", getUserWithdrawals(logger, services.UsersService))
		})
	})

	return r
}

func notAllowed() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := models.NewErrNotAllowed(
			errors.New("method not allowed"),
			"HTTP_ERROR",
			"Method not allowed",
		)
		_ = models.WriteStatusError(w, err)
	}
}

func notFound() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := models.NewErrNotFound(
			errors.New("endpoint not found"),
			"HTTP_ERROR",
			"Endpoint not found",
		)
		_ = models.WriteStatusError(w, err)
	}
}

func liveProbe(logger zerolog.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("."))
		if err != nil {
			logger.Error().Err(err).Send()
		}
	}
}

func readyProbe(logger zerolog.Logger, health services.HealthService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		if err := health.ConnPing(r.Context()); err != nil {
			logger.Error().Err(err).Send()
			err = models.WriteStatusError(w, err)
			logger.Error().Err(err).Send()
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("."))
		if err != nil {
			logger.Error().Err(err).Send()
		}
	}
}
