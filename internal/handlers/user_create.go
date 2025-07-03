package handlers

import (
	"encoding/json"
	"go-crud/internal/domain"
	"net/http"
)

type UserCreateHandler struct {
	UserRepository domain.UserRepository
}

func NewUserCreateHandler(userRepository domain.UserRepository) *UserCreateHandler {
	return &UserCreateHandler{
		UserRepository: userRepository,
	}
}

func (h *UserCreateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err := h.UserRepository.Create(&user)
	if err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}
