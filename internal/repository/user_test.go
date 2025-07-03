package repository

import (
	"go-crud/internal/domain"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}

	defer db.Close()

	repo := NewUserRepository(db)

	subtests := []struct {
		name        string
		user        *domain.User
		expectedID  int
		expectedErr error
	}{
		{
			name: "success",
			user: &domain.User{
				Username: "testuser",
				Email:    "test@email.com",
				Password: "hashedpassword123",
			},
			expectedID:  1,
			expectedErr: nil,
		},
		// TODO: duplicate record?
	}
	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			query := mock.ExpectQuery(`INSERT INTO users`).
				WithArgs(subtest.user.Username, subtest.user.Email, subtest.user.Password)

			switch subtest.name {
			case "success":
				fixedTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
					AddRow(subtest.expectedID, fixedTime, fixedTime)
				query.WillReturnRows(rows)

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
			}
		})
	}
}
