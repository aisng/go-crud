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

type MethodHandlers map[string]http.HandlerFunc

func NewHandler(deps Dependencies) *Handler {
	return &Handler{
		User: NewUserHandler(deps.UserRepo),
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/users", MethodRouter(MethodHandlers{http.MethodPost: h.User.Create}))
	mux.HandleFunc("/users/{id}", MethodRouter(MethodHandlers{
		http.MethodGet:    h.User.GetByID,
		http.MethodPut:    h.User.Update,
		http.MethodDelete: h.User.Delete,
	}))
}

func MethodRouter(handlers MethodHandlers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if h, ok := handlers[r.Method]; ok {
			h(w, r)
			return
		}
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func parsePathParam(r *http.Request, base string) (string, error) {
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	if len(parts) != 2 || parts[0] != base {
		return "", fmt.Errorf("invalid path format")
	}

	return parts[1], nil
}
