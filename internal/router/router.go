// Package router serves mux with registered routes
package router

import (
	"go-crud/internal/handler"
	"net/http"
)

func NewRouter(h *handler.Handler) *http.ServeMux {
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	return mux
}
