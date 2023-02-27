package handlers

import (
	"io"
	"net/http"

	"github.com/c0dered273/go-musthave-diploma-tpl/internal/models"
	"github.com/c0dered273/go-musthave-diploma-tpl/internal/services"
	"github.com/rs/zerolog"
)

func addOrders(logger zerolog.Logger, service services.UsersService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		if contentType != "text/plain" {
			logger.Error().Err(ErrContentType).Send()
			models.WriteStatusError(w, ErrContentType)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error().Err(err).Send()
			models.WriteStatusError(w, ErrParseRequest)
			return
		}
		defer r.Body.Close()

		err = service.CreateOrders(r.Context(), string(body))
		if err != nil {
			models.WriteStatusError(w, err)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}
