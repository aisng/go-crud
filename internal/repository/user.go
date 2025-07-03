package repository

import (
	"database/sql"
	"go-crud/internal/model"
)

type UserRpository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRpository {
	return &UserRpository{db: db}
}

func Create(user *model.User) error {
	return nil
}
