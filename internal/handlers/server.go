package handlers

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/c0dered273/go-musthave-diploma-tpl/internal/configs"
)

const (
	readHeaderTimeout = 5 * time.Second
	readTimeOut       = 1 * time.Minute
	writeTimeout      = 1 * time.Minute
	idleTimeout       = 2 * time.Minute
)

func NewServer(ctx context.Context, cfg *configs.ServerConfig, handler http.Handler) *http.Server {
	return &http.Server{
		Addr: cfg.RunAddress,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
		Handler:           handler,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeOut,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}
}
