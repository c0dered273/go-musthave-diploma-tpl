package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/c0dered273/go-musthave-diploma-tpl/internal/models"
	"github.com/c0dered273/go-musthave-diploma-tpl/internal/services"
	"github.com/mailru/easyjson"
	"github.com/rs/zerolog"
)

var (
	ErrParseRequest = models.NewErrBadRequest(nil, "BAD_REQUEST", "Failed to parse request")
	ErrNoAuthHeader = models.NewErrBadRequest(nil, "BAD_REQUEST", "No authorization header")
	ErrServerError  = models.NewErrInternal(nil, "SERVER_ERROR", "Internal error")
)

// TODO("Убрать копипасту")

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

func withdrawals(logger zerolog.Logger, service services.UsersService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if len(authorization) == 0 {
			models.WriteStatusError(w, ErrNoAuthHeader)
			return
		}

		tokenString := strings.Split(authorization, "Bearer ")[1]

		// TODO("Implement")

		usrname, err := service.GetWithdrawals(r.Context(), tokenString)
		if err != nil {
			models.WriteStatusError(w, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(usrname))
		if err != nil {
			logger.Error().Err(err).Send()
			return
		}
	}
}
