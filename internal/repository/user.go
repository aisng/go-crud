package repository

import (
	"database/sql"
	"go-crud/internal/domain"
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
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = int(id)

	row := r.db.QueryRow("SELECT created_at, updated_at FROM users WHERE id = ?", user.ID)
	err = row.Scan(&user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetByID(id int) (*domain.User, error) {
	return nil, nil
}

func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	return nil, nil
}

func (r *UserRepository) Update(user *domain.User) error {
	return nil
}

func (r *UserRepository) Delete(id int) error {
	return nil
}
