// Package handler handles HTTP requests and registers routes
package handler

import (
	"fmt"
	"go-crud/internal/domain"
	"net/http"
	"strings"
)

type ResouceHandler interface {
	GetByID(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

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
	mux.HandleFunc("/users/{id}", resourceByIDHandler(h.User))
}

func resourceByIDHandler(handler ResouceHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetByID(w, r)
		case http.MethodPut:
			handler.Update(w, r)
		case http.MethodDelete:
			handler.Delete(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
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
