package domain

import "time"

type User struct {
	ID        int64     `json:"id"         db:"id"`
	Email     string    `json:"email"      db:"email"`
	Username  string    `json:"username"   db:"username"`
	Password  string    `json:"-"          db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type UserRepository interface {
	Create(user *User) error
	GetByID(id int64) (*User, error)
	Update(id int64, upd *UserUpdate) error
	Delete(id int64) error
}

type UserUpdate struct {
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email,omitempty"`
	Password *string `json:"password,omitempty"`
}
