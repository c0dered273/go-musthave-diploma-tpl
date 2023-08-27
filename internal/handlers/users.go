package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/c0dered273/go-musthave-diploma-tpl/internal/models"
	"github.com/c0dered273/go-musthave-diploma-tpl/internal/services"
	"github.com/mailru/easyjson"
	"github.com/rs/zerolog"
)

var (
	ErrParseRequest = models.NewErrBadRequest(nil, "BAD_REQUEST", "Failed to parse request")
	ErrContentType  = models.NewErrBadRequest(nil, "BAD_REQUEST", "Wrong content type")
	ErrServerError  = models.NewErrInternal(nil, "SERVER_ERROR", "Internal error")
)

func registerUser(logger zerolog.Logger, service services.UsersService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error().Err(err).Send()
			err = models.WriteStatusError(w, ErrParseRequest)
			logger.Error().Err(err).Send()
			return
		}
		defer r.Body.Close()

		newUser := &models.LoginRequestDTO{}
		err = easyjson.Unmarshal(body, newUser)
		if err != nil {
			logger.Error().Err(err).Send()
			err = models.WriteStatusError(w, ErrParseRequest)
			logger.Error().Err(err).Send()
			return
		}

		authResponse, err := service.NewUser(r.Context(), newUser)
		if err != nil {
			err = models.WriteStatusError(w, err)
			logger.Error().Err(err).Send()
			return
		}
		body, err = easyjson.Marshal(authResponse)
		if err != nil {
			logger.Error().Err(err).Send()
			err = models.WriteStatusError(w, ErrServerError)
			logger.Error().Err(err).Send()
			return
		}

		w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", authResponse.AccessToken))
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(body)
		if err != nil {
			logger.Error().Err(err).Send()
			return
		}
	}
}

func loginUser(logger zerolog.Logger, service services.UsersService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error().Err(err).Send()
			err = models.WriteStatusError(w, ErrParseRequest)
			logger.Error().Err(err).Send()
			return
		}
		defer r.Body.Close()

		newUser := &models.LoginRequestDTO{}
		err = easyjson.Unmarshal(body, newUser)
		if err != nil {
			logger.Error().Err(err).Send()
			err = models.WriteStatusError(w, ErrParseRequest)
			logger.Error().Err(err).Send()
			return
		}

		authResponse, err := service.LoginUser(r.Context(), newUser)
		if err != nil {
			err = models.WriteStatusError(w, err)
			logger.Error().Err(err).Send()
			return
		}
		body, err = easyjson.Marshal(authResponse)
		if err != nil {
			logger.Error().Err(err).Send()
			err = models.WriteStatusError(w, ErrServerError)
			logger.Error().Err(err).Send()
			return
		}

		w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", authResponse.AccessToken))
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(body)
		if err != nil {
			logger.Error().Err(err).Send()
			return
		}
	}
}

func getUserOrders(logger zerolog.Logger, service services.UsersService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		orders, err := service.GetOrders(r.Context())
		if err != nil {
			err = models.WriteStatusError(w, err)
			logger.Error().Err(err).Send()
			return
		}

		if len(orders) == 0 {
			http.Error(w, "no content", http.StatusNoContent)
			return
		}

		ordersResponse, err := easyjson.Marshal(orders)
		if err != nil {
			logger.Error().Err(err).Send()
			err = models.WriteStatusError(w, ErrServerError)
			logger.Error().Err(err).Send()
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(ordersResponse)
		if err != nil {
			logger.Error().Err(err).Send()
			return
		}
	}
}

func getUserWithdrawals(logger zerolog.Logger, service services.UsersService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		withdrawals, err := service.GetWithdrawals(r.Context())
		if err != nil {
			err = models.WriteStatusError(w, err)
			logger.Error().Err(err).Send()
			return
		}

		if len(withdrawals) == 0 {
			http.Error(w, "no content", http.StatusNoContent)
			return
		}

		withdrawalsResponse, err := easyjson.Marshal(withdrawals)
		if err != nil {
			logger.Error().Err(err).Send()
			err = models.WriteStatusError(w, ErrServerError)
			logger.Error().Err(err).Send()
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(withdrawalsResponse)
		if err != nil {
			logger.Error().Err(err).Send()
			return
		}
	}
}

func getUserBalance(logger zerolog.Logger, service services.UsersService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		balance, err := service.GetBalance(r.Context())
		if err != nil {
			err = models.WriteStatusError(w, err)
			logger.Error().Err(err).Send()
			return
		}

		balanceResponse, err := easyjson.Marshal(balance)
		if err != nil {
			logger.Error().Err(err).Send()
			err = models.WriteStatusError(w, ErrServerError)
			logger.Error().Err(err).Send()
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write(balanceResponse)
		if err != nil {
			logger.Error().Err(err).Send()
			return
		}
	}
}
