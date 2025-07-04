package repository

import (
	"fmt"
	"go-crud/internal/domain"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestUserRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}

	defer db.Close()

	repo := NewUserRepository(db)
	fixedTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	user := &domain.User{
		Username: "testuser",
		Email:    "test@email.com",
		Password: "hashedpassword123",
	}

	subtests := []struct {
		name        string
		user        *domain.User
		expectedID  int64
		expectedErr error
		setupMock   func()
	}{
		{
			name:        "Create success",
			user:        user,
			expectedID:  1,
			expectedErr: nil,
			setupMock: func() {
				mock.ExpectExec(`INSERT INTO users`).
					WithArgs(user.Username, user.Email, user.Password).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectQuery(`SELECT created_at, updated_at FROM users WHERE id = ?`).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at"}).
						AddRow(fixedTime, fixedTime))
			},
		},
		{
			name:        "Create duplicate entry",
			user:        user,
			expectedID:  0,
			expectedErr: fmt.Errorf("Duplicate entry"),
			setupMock: func() {
				mock.ExpectExec(`INSERT INTO users`).
					WithArgs(user.Username, user.Email, user.Password).
					WillReturnError(fmt.Errorf("Duplicate entry"))
			},
		},
	}
	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			switch subtest.name {
			case "success":
				subtest.setupMock()

				err := repo.Create(subtest.user)
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
				if subtest.user.ID != subtest.expectedID {
					t.Errorf("Expected ID: %d, got: %d", subtest.expectedID, subtest.user.ID)
				}
				if subtest.user.CreatedAt.IsZero() {
					t.Error("Expected CreatedAt to be set")
				}
				if subtest.user.UpdatedAt.IsZero() {
					t.Error("Expected UpdatedAt to be set")
				}
			case "duplicate entry":
				err := repo.Create(subtest.user)

				if err != subtest.expectedErr {
					t.Errorf("expected duplicate entry error, got: %v", err)
				}

			}
		})
	}
}
func strPtr(s string) *string { return &s }

func TestUserRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)
	fixedTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	subtests := []struct {
		name         string
		id           int64
		expectedUser *domain.User
		expectedErr  error
		setupMock    func()
	}{
		{
			name: "user found",
			id:   1,
			expectedUser: &domain.User{
				ID:        1,
				Username:  "testuser",
				Email:     "test@email.com",
				CreatedAt: fixedTime,
				UpdatedAt: fixedTime,
			},
			expectedErr: nil,
			setupMock: func() {
				rows := sqlmock.NewRows([]string{"id", "username", "email", "created_at", "updated_at"}).
					AddRow(1, "testuser", "test@email.com", fixedTime, fixedTime)

				mock.ExpectQuery(`SELECT id, username, email, created_at, updated_at FROM users WHERE id = \?`).
					WithArgs(1).
					WillReturnRows(rows)
			},
		},
		{
			name:         "user not found",
			id:           9999,
			expectedUser: nil,
			expectedErr:  fmt.Errorf("user not found"),
			setupMock: func() {
				mock.ExpectQuery(`SELECT id, username, email, created_at, updated_at FROM users WHERE id = \?`).
					WithArgs(9999).
					WillReturnError(fmt.Errorf("user not found"))
			},
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			subtest.setupMock()

			user, err := repo.GetByID(subtest.id)

			if (subtest.expectedErr == nil && err != nil) || (subtest.expectedErr != nil && err == nil) {
				t.Errorf("expected error: %v, got: %v", subtest.expectedErr, err)
			}

			if !reflect.DeepEqual(subtest.expectedUser, user) {
				t.Errorf("expected user: %v, got: %v", subtest.expectedUser, user)
			}
		})
	}
}

func TestUserRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	repo := NewUserRepository(db)

	subtests := []struct {
		name        string
		id          int64
		update      *domain.UserUpdate
		expectedErr error
		setupMock   func()
	}{
		{
			name: "user found and updated",
			id:   1,
			update: &domain.UserUpdate{
				Email:    strPtr("new@email.com"),
				Username: strPtr("newuser"),
			},
			setupMock: func() {
				mock.ExpectExec(`UPDATE users SET username = \?, email = \?, updated_at = NOW\(\) WHERE id = \?`).
					WithArgs("newuser", "new@email.com", 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedErr: nil,
		},
		{
			name: "user not found",
			id:   999,
			update: &domain.UserUpdate{
				Email: strPtr("notfound@email.com"),
			},
			setupMock: func() {
				mock.ExpectExec(`UPDATE users SET email = \?, updated_at = NOW\(\) WHERE id = \?`).
					WithArgs("notfound@email.com", 999).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedErr: fmt.Errorf("user not found"),
		},
		{
			name:   "no fields to update",
			id:     2,
			update: &domain.UserUpdate{},
			setupMock: func() {
			},
			expectedErr: nil,
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			subtest.setupMock()

			err := repo.Update(subtest.id, subtest.update)

			if (subtest.expectedErr == nil && err != nil) || (subtest.expectedErr != nil && err == nil) {
				t.Errorf("expected error: %v, got: %v", subtest.expectedErr, err)
			}
		})
	}
}
