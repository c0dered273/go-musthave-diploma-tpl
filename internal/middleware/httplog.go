package middleware

import (
	"net/http"

	"github.com/rs/zerolog"
)

// TODO("Реализовать нормальный логер для http на основе github.com/go-chi/httplog")

func HTTPLog(logger zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {

			}()

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
