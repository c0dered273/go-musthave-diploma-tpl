package repositories

import (
	"errors"
	"strings"
)

var (
	ErrAlreadyExists = errors.New("already exists")
	ErrNotFound      = errors.New("not found")
)

func strip(str string) string {
	rsl := strings.ReplaceAll(str, "\t", "")
	return rsl
}
