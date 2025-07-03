// Package handler handles HTTP requests and registers routes
package handler

import (
	"fmt"
	"go-crud/internal/domain"
	"net/http"
	"strings"
)

type Dependencies struct {
	UserRepo domain.UserRepository
}

type Handler struct {
	User *UserHandler
}

func NewHandler(deps Dependencies) *Handler {
	return &Handler{
		User: NewUserHandler(deps.UserRepo),
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/users", h.User.Create)
	mux.HandleFunc("/users/{id}", h.User.GetByID)
}

func parsePathParam(r *http.Request, base string) (string, error) {
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	if len(parts) != 2 || parts[0] != base {
		return "", fmt.Errorf("invalid path format")
	}

	return parts[1], nil
}
