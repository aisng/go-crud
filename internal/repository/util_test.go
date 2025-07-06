package repository

import (
	"database/sql"
	"fmt"
	"go-crud/internal/domain"
	"testing"

	"github.com/go-sql-driver/mysql"
)

func TestResolveSQLError(t *testing.T) {
	tests := []struct {
		name     string
		inputErr error
		expected error
	}{
		{
			name:     "nil error",
			inputErr: nil,
			expected: nil,
		},
		{
			name:     "sql.ErrNoRows",
			inputErr: sql.ErrNoRows,
			expected: domain.ErrNotFound,
		},
		{
			name:     "mysql duplicate entry error",
			inputErr: &mysql.MySQLError{Number: 1062, Message: "Duplicate entry"},
			expected: domain.ErrAlreadyExists,
		},
		{
			name:     "generic database error",
			inputErr: fmt.Errorf("connection lost"),
			expected: fmt.Errorf("unexpected db error: connection lost"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := resolveSQLError(test.inputErr)

			if test.expected == nil && result == nil {
				return
			}

			if (test.expected == nil) != (result == nil) {
				t.Errorf("expected error: %v, got: %v", test.expected, result)
				return
			}

			if test.expected.Error() != result.Error() {
				t.Errorf("expected error: %v, got: %v", test.expected, result)
			}
		})
	}
}
