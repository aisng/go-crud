package handler

import (
	"encoding/json"
	"fmt"
	"go-crud/internal/domain"
	"net/http"
	"strconv"
)

type UserHandler struct {
	userRepo domain.UserRepository
}

func NewUserHandler(userRepository domain.UserRepository) *UserHandler {
	return &UserHandler{
		userRepo: userRepository,
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err := h.userRepo.Create(&user)
	if err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr, err := parsePathParam(r, "users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid parameter 'id'", http.StatusBadRequest)
		return
	}

	user, err := h.userRepo.GetByID(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	// if r.Method != http.MethodPut {
	// 	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// 	return
	// }

	idStr, err := parsePathParam(r, "users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid parameter 'id", http.StatusBadRequest)
		return
	}

	var userUpd domain.UserUpdate
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&userUpd); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = h.userRepo.Update(id, &userUpd)
	if err != nil {
		if err.Error() == "user not found" {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to update user with id %d, err: %v", id, err), http.StatusInternalServerError)
		return
	}

	user, err := h.userRepo.GetByID(id)
	if err != nil {
		http.Error(w, fmt.Sprintf("User with id %d updated but not retrieved", id), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr, err := parsePathParam(r, "users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid parameter 'id'", http.StatusBadRequest)
		return
	}

	err = h.userRepo.Delete(id)
	if err != nil {
		if err.Error() == "user not found" {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
