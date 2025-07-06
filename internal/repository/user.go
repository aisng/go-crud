// Package repository
package repository

import (
	"database/sql"
	"go-crud/internal/domain"
	"strings"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User) error {
	query := `
	INSERT INTO users (username, email, password, created_at, updated_at)
	VALUES (?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.Exec(query, user.Username, user.Email, user.Password)
	if err != nil {
		return resolveSQLError(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return resolveSQLError(err)
	}

	user.ID = id

	row := r.db.QueryRow("SELECT created_at, updated_at FROM users WHERE id = ?", user.ID)
	err = row.Scan(&user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return resolveSQLError(err)
	}
	return nil
}

func (r *UserRepository) GetByID(id int64) (*domain.User, error) {
	query := `
	SELECT id, username, email, created_at, updated_at FROM users
	WHERE id = ?`

	row := r.db.QueryRow(query, id)

	var user domain.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, resolveSQLError(err)
	}

	return &user, nil
}

func (r *UserRepository) Update(id int64, upd *domain.UserUpdate) error {
	setClauses := []string{}
	args := []any{}

	if upd.Username != nil {
		setClauses = append(setClauses, "username = ?")
		args = append(args, *upd.Username)
	}

	if upd.Email != nil {
		setClauses = append(setClauses, "email = ?")
		args = append(args, *upd.Email)
	}

	if upd.Password != nil {
		setClauses = append(setClauses, "password = ?")
		args = append(args, *upd.Password)
	}

	if len(setClauses) == 0 {
		return nil
	}

	setClauses = append(setClauses, "updated_at = NOW()")
	query := "UPDATE users SET " + strings.Join(setClauses, ", ") + " WHERE id = ?"
	args = append(args, id)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return resolveSQLError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return resolveSQLError(err)
	}

	if rowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *UserRepository) Delete(id int64) error {
	query := "DELETE FROM users WHERE id = ?"

	result, err := r.db.Exec(query, id)
	if err != nil {
		return resolveSQLError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return resolveSQLError(err)
	}

	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}
