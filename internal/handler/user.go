package handler

import (
	"encoding/json"
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
		WriteError(w, ErrInvalidJSON.Message, ErrInvalidJSON.Code)
		return
	}

	err := h.userRepo.Create(&user)
	if err != nil {
		handleDomainError(w, err)
		return
	}

	WriteResponse(w, user, http.StatusCreated)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr, err := parsePathParam(r, "users")
	if err != nil {
		WriteError(w, ErrInvalidPath.Message, ErrInvalidPath.Code)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		WriteError(w, ErrInvalidID.Message, ErrInvalidID.Code)
		return
	}

	user, err := h.userRepo.GetByID(id)
	if err != nil {
		handleDomainError(w, err)
		return
	}

	WriteResponse(w, user, http.StatusOK)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr, err := parsePathParam(r, "users")
	if err != nil {
		WriteError(w, ErrInvalidPath.Message, ErrInvalidPath.Code)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		WriteError(w, ErrInvalidID.Message, ErrInvalidID.Code)
		return
	}

	var userUpd domain.UserUpdate
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&userUpd); err != nil {
		WriteError(w, ErrInvalidJSON.Message, ErrInvalidJSON.Code)
		return
	}

	err = h.userRepo.Update(id, &userUpd)
	if err != nil {
		handleDomainError(w, err)
		return
	}

	user, err := h.userRepo.GetByID(id)
	if err != nil {
		handleDomainError(w, err)
		return
	}

	WriteResponse(w, user, http.StatusOK)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr, err := parsePathParam(r, "users")
	if err != nil {
		WriteError(w, ErrInvalidPath.Message, ErrInvalidPath.Code)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		WriteError(w, ErrInvalidID.Message, ErrInvalidID.Code)
		return
	}

	err = h.userRepo.Delete(id)
	if err != nil {
		handleDomainError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
