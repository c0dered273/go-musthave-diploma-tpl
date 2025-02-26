package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/c0dered273/go-musthave-diploma-tpl/internal/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/rs/zerolog"
)

var (
	ClaimCtxKey CtxKeyType = "claim"
)

type CtxKeyType string

var (
	ErrInvalidToken = models.NewErrUnauthorized(nil, "AUTH_ERROR", "Access token invalid")
)

// TODO("To refactor")

func JwtVerifier(logger zerolog.Logger, secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			tokenString, err := tokenFromHeader(r)
			if err != nil {
				err = models.WriteStatusError(w, err)
				logger.Error().Err(err).Send()
				return
			}

			claims, err := validateToken(tokenString, secret)
			if err != nil {
				logger.Error().Err(err).Send()
				err = models.WriteStatusError(w, ErrInvalidToken)
				logger.Error().Err(err).Send()
				return
			}

			ctx := context.WithValue(r.Context(), ClaimCtxKey, claims)

			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func tokenFromHeader(r *http.Request) (string, error) {
	bearer := r.Header.Get("Authorization")
	if len(bearer) > 7 && strings.ToUpper(bearer[0:6]) == "BEARER" {
		return bearer[7:], nil
	}
	return "", ErrInvalidToken
}

func validateToken(tokenString string, secret string) (*models.AuthClaim, error) {
	claims := &models.AuthClaim{RegisteredClaims: jwt.RegisteredClaims{}}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if _, ok := token.Claims.(*models.AuthClaim); !(ok && token.Valid) {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
