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
			models.WriteStatusError(w, ErrParseRequest)
			return
		}
		defer r.Body.Close()

		newUser := &models.LoginRequestDTO{}
		err = easyjson.Unmarshal(body, newUser)
		if err != nil {
			logger.Error().Err(err).Send()
			models.WriteStatusError(w, ErrParseRequest)
			return
		}

		authResponse, err := service.NewUser(r.Context(), newUser)
		if err != nil {
			models.WriteStatusError(w, err)
			return
		}
		body, err = easyjson.Marshal(authResponse)
		if err != nil {
			logger.Error().Err(err).Send()
			models.WriteStatusError(w, ErrServerError)
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
			models.WriteStatusError(w, ErrParseRequest)
			return
		}
		defer r.Body.Close()

		newUser := &models.LoginRequestDTO{}
		err = easyjson.Unmarshal(body, newUser)
		if err != nil {
			logger.Error().Err(err).Send()
			models.WriteStatusError(w, ErrParseRequest)
			return
		}

		authResponse, err := service.LoginUser(r.Context(), newUser)
		if err != nil {
			models.WriteStatusError(w, err)
			return
		}
		body, err = easyjson.Marshal(authResponse)
		if err != nil {
			logger.Error().Err(err).Send()
			models.WriteStatusError(w, ErrServerError)
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

		w.WriteHeader(http.StatusOK)
	}
}

func getUserOrders(logger zerolog.Logger, service services.UsersService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		orders, err := service.GetUserOrders(r.Context())
		if err != nil {
			models.WriteStatusError(w, err)
			return
		}

		if len(orders) == 0 {
			http.Error(w, "no content", http.StatusNoContent)
			return
		}

		ordersResponse, err := easyjson.Marshal(orders)

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
			models.WriteStatusError(w, err)
			return
		}

		if len(withdrawals) == 0 {
			http.Error(w, "no content", http.StatusNoContent)
		}

		withdrawalsResponse, err := easyjson.Marshal(withdrawals)
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(withdrawalsResponse)
		if err != nil {
			logger.Error().Err(err).Send()
			return
		}
	}
}
