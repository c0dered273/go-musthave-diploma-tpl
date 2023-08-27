package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/c0dered273/go-musthave-diploma-tpl/internal/clients"
	"github.com/c0dered273/go-musthave-diploma-tpl/internal/configs"
	"github.com/c0dered273/go-musthave-diploma-tpl/internal/handlers"
	"github.com/c0dered273/go-musthave-diploma-tpl/internal/loggers"
	"github.com/c0dered273/go-musthave-diploma-tpl/internal/repositories"
	"github.com/c0dered273/go-musthave-diploma-tpl/internal/services"
	"github.com/c0dered273/go-musthave-diploma-tpl/internal/store"
	"github.com/c0dered273/go-musthave-diploma-tpl/internal/validators"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
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

var (
	logger zerolog.Logger
	conn   *pgxpool.Pool
	server *http.Server
)

func main() {
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	serverInitAndStart(serverCtx)

	// graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-shutdown
		shutdownCtx, shutdownCancelCtx := context.WithTimeout(serverCtx, shutdownTimeout)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				logger.Fatal().Msg("server: graceful shutdown timed out.. forcing exit.")
			}
		}()

		logger.Info().Msg("server: shutting down..")
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			logger.Error().Err(err).Msg("server: graceful shutdown failed")
		}

		conn.Close()

		serverStopCtx()
		shutdownCancelCtx()
	}()

	<-serverCtx.Done()
}

func serverInitAndStart(ctx context.Context) {
	// logger, validator, config
	initLogger := loggers.NewDefaultLogger()
	initLogger.Info().Msg("server: init")
	validator := validators.NewValidatorTagName("mapstructure")
	cfg, err := configs.NewServerConfig(configFileName, configFilePath, initLogger, validator)
	if err != nil {
		initLogger.Fatal().Err(err).Msg("server: config init failed")
	}
	logger = loggers.NewServerLogger(cfg)

	//migration
	err = repositories.ApplyMigration(logger, cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("server: DB migration init failed")
	}

	//REST client
	accrualClient := clients.NewAccrualClient(cfg)

	// repositories
	conn, err = store.NewPgxConn(ctx, logger, cfg)
	connCheck := store.NewPgxConnCheck(conn)
	if err != nil {
		logger.Fatal().Err(err).Msg("server: DB connection init failed")
	}
	usersRepo := repositories.NewUserRepository(conn)
	ordersRepo := repositories.NewOrderRepository(conn)
	withdrawalsRepo := repositories.NewWithdrawalRepository(conn)

	// services
	serviceContext := &services.ServiceContext{
		HealthService: services.NewHealthService(logger, connCheck),
		UsersService:  services.NewUsersService(logger, cfg, validator, usersRepo, ordersRepo, withdrawalsRepo, accrualClient),
	}

	// http server
	handler := handlers.NewHandler(logger, cfg, serviceContext)
	server = handlers.NewServer(ctx, cfg, handler)
	ln, err := net.Listen("tcp", cfg.RunAddress)
	if err != nil {
		logger.Fatal().Err(err).Msgf("server: failed to start server on %s", cfg.RunAddress)
	}
	logger.Info().Msgf("server: listening %s", cfg.RunAddress)

	go func() {
		err = server.Serve(ln)
		if err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				logger.Fatal().Err(err).Msgf("server: failed to start server on %s", cfg.RunAddress)
			}
		}
	}()
}
