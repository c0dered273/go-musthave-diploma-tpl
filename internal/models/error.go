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
	Timestamp time.Time `json:"timestamp,omitempty"`
	ErrorCode string    `json:"errorCode,omitempty"`
	Message   string    `json:"message,omitempty"`
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

func NewStatusError(err error, httpCode int, errorCode string, message string) *StatusError {
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
		StatusError: NewStatusError(
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
		StatusError: NewStatusError(
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
		StatusError: NewStatusError(
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
		StatusError: NewStatusError(
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
		StatusError: NewStatusError(
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
		StatusError: NewStatusError(
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
		StatusError: NewStatusError(
			err,
			http.StatusConflict,
			errorCode,
			message,
		),
	}
}

type HttpStatusCreated struct {
	*StatusError
}

func NewStatusCreated(message string) *HttpStatusCreated {
	return &HttpStatusCreated{
		StatusError: NewStatusError(
			nil,
			http.StatusAccepted,
			"",
			message,
		),
	}
}