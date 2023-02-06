package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/c0dered273/go-musthave-diploma-tpl/internal/configs"
	"github.com/c0dered273/go-musthave-diploma-tpl/internal/handlers"
	"github.com/c0dered273/go-musthave-diploma-tpl/internal/loggers"
	"github.com/c0dered273/go-musthave-diploma-tpl/internal/repositories"
	"github.com/c0dered273/go-musthave-diploma-tpl/internal/services"
	"github.com/c0dered273/go-musthave-diploma-tpl/internal/validators"
)

const (
	shutdownTimeout = 20 * time.Second
)

var (
	configFileName = "application"
	configFilePath = []string{
		".",
		"./configs/",
	}
)

func main() {
	serverCtx, serverStopCtx := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	// logger, validator, config
	logger := loggers.NewDefaultLogger()
	logger.Info().Msg("server: init")
	validator := validators.NewValidatorTagName("mapstructure")
	cfg, err := configs.NewServerConfig(configFileName, configFilePath, logger, validator)
	if err != nil {
		logger.Fatal().Err(err).Msg("server: config init failed")
	}
	srvLogger := loggers.NewServerLogger(cfg)

	// repository
	repo, repErr := repositories.NewCrudRepository(serverCtx, logger, cfg)
	if repErr != nil {
		logger.Fatal().Err(repErr).Msg("server: DB connection init failed")
	}

	// services
	serviceContext := services.ServiceContext{
		HealthService: services.NewHealthService(repo, srvLogger),
	}

	// http server
	handler := handlers.NewHandler(srvLogger, serviceContext)
	server := handlers.NewServer(serverCtx, cfg, handler)
	ln, err := net.Listen("tcp", cfg.RunAddress)
	if err != nil {
		logger.Fatal().Err(err).Msgf("server: failed to start server on %s", cfg.RunAddress)
	}
	srvLogger.Info().Msgf("server: listening %s", cfg.RunAddress)

	go func() {
		err = server.Serve(ln)
		if err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				logger.Fatal().Err(err).Msgf("server: failed to start server on %s", cfg.RunAddress)
			}
		}
	}()

	// graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-shutdown
		shutdownCtx, shutdownCancelCtx := context.WithTimeout(serverCtx, shutdownTimeout)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				srvLogger.Fatal().Msg("server: graceful shutdown timed out.. forcing exit.")
			}
		}()

		srvLogger.Info().Msg("server: shutting down..")
		repo.Close()
		err = server.Shutdown(shutdownCtx)
		if err != nil {
			srvLogger.Fatal().Err(err).Msg("server: graceful shutdown failed")
		}

		serverStopCtx()
		shutdownCancelCtx()
	}()

	<-serverCtx.Done()
	wg.Wait()
}
