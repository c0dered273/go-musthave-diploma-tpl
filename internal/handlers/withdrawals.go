package handlers

import (
	"io"
	"net/http"

	"github.com/c0dered273/go-musthave-diploma-tpl/internal/models"
	"github.com/c0dered273/go-musthave-diploma-tpl/internal/services"
	"github.com/mailru/easyjson"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"
)

func withdrawBalance(logger zerolog.Logger, service services.UsersService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error().Err(err).Send()
			err = models.WriteStatusError(w, ErrParseRequest)
			logger.Error().Err(err).Send()
			return
		}
		defer r.Body.Close()

		request := models.WithdrawRequest{}
		err = easyjson.Unmarshal(body, &request)
		if err != nil {
			logger.Error().Err(err).Send()
			err = models.WriteStatusError(w, ErrParseRequest)
			logger.Error().Err(err).Send()
			return
		}

		err = service.WithdrawBalance(r.Context(), request.OrderID, decimal.NewFromFloat(request.Sum))
		if err != nil {
			err = models.WriteStatusError(w, err)
			logger.Error().Err(err).Send()
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
