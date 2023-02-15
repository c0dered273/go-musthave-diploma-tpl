package models

import (
	"errors"
	"net/http"
	"time"

	"github.com/mailru/easyjson"
	_ "github.com/mailru/easyjson/gen"
)

type HttpError interface {
	HttpError(w http.ResponseWriter) error
}

//go:generate easyjson error.go
//easyjson:json
type StatusError struct {
	Err       error     `json:"-"`
	HttpCode  int       `json:"-"`
	Timestamp time.Time `json:"timestamp"`
	ErrorCode string    `json:"errorCode"`
	Message   string    `json:"message"`
}

func (se *StatusError) Error() string {
	return se.Err.Error()
}

func (se *StatusError) HttpError(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(se.HttpCode)
	_, _, err := easyjson.MarshalToHTTPResponseWriter(se, w)
	if err != nil {
		return err
	}
	return nil
}

func newStatusError(err error, httpCode int, errorCode string, message string) *StatusError {
	return &StatusError{
		Err:       err,
		HttpCode:  httpCode,
		Timestamp: time.Now(),
		ErrorCode: errorCode,
		Message:   message,
	}
}

func WriteStatusError(w http.ResponseWriter, err error) {
	if statusErr, ok := err.(HttpError); ok {
		errSts := statusErr.HttpError(w)
		if errSts != nil {
			panic(err)
		}
	} else {
		panic(errors.New("status error: failed to cast err to HttpError"))
	}
}

type ErrInternal struct {
	*StatusError
}

func NewErrInternal(err error, errorCode string, message string) *ErrInternal {
	return &ErrInternal{
		StatusError: newStatusError(
			err,
			http.StatusInternalServerError,
			errorCode,
			message,
		),
	}
}

type ErrBadRequest struct {
	*StatusError
}

func NewErrBadRequest(err error, errorCode string, message string) *ErrBadRequest {
	return &ErrBadRequest{
		StatusError: newStatusError(
			err,
			http.StatusBadRequest,
			errorCode,
			message,
		),
	}
}

type ErrUnauthorized struct {
	*StatusError
}

func NewErrUnauthorized(err error, errorCode string, message string) *ErrUnauthorized {
	return &ErrUnauthorized{
		StatusError: newStatusError(
			err,
			http.StatusUnauthorized,
			errorCode,
			message,
		),
	}
}

type ErrForbidden struct {
	*StatusError
}

func NewErrForbidden(err error, errorCode string, message string) *ErrForbidden {
	return &ErrForbidden{
		StatusError: newStatusError(
			err,
			http.StatusForbidden,
			errorCode,
			message,
		),
	}
}

type ErrNotFound struct {
	*StatusError
}

func NewErrNotFound(err error, errorCode string, message string) *ErrNotFound {
	return &ErrNotFound{
		StatusError: newStatusError(
			err,
			http.StatusNotFound,
			errorCode,
			message,
		),
	}
}

type ErrNotAllowed struct {
	*StatusError
}

func NewErrNotAllowed(err error, errorCode string, message string) *ErrNotAllowed {
	return &ErrNotAllowed{
		StatusError: newStatusError(
			err,
			http.StatusMethodNotAllowed,
			errorCode,
			message,
		),
	}
}

type ErrConflict struct {
	*StatusError
}

func NewErrConflict(err error, errorCode string, message string) *ErrConflict {
	return &ErrConflict{
		StatusError: newStatusError(
			err,
			http.StatusConflict,
			errorCode,
			message,
		),
	}
}
