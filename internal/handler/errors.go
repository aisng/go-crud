package handler

import (
	"errors"
	"go-crud/internal/domain"
	"net/http"
)

type HTTPError struct {
	Message string
	Code    int
}

func (e *HTTPError) Error() string {
	return e.Message
}

var (
	ErrInvalidJSON = &HTTPError{Message: "invalid json", Code: http.StatusBadRequest}
	ErrInvalidID   = &HTTPError{Message: "invalid parameter 'id'", Code: http.StatusBadRequest}
	ErrInvalidPath = &HTTPError{Message: "invalid path format", Code: http.StatusBadRequest}
)

func handleDomainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		WriteError(w, domain.ErrNotFound.Error(), http.StatusNotFound)
	case errors.Is(err, domain.ErrAlreadyExists):
		WriteError(w, domain.ErrAlreadyExists.Error(), http.StatusConflict)
	default:
		WriteError(w, "internal server error", http.StatusInternalServerError)
	}
}
