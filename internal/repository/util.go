package repository

import (
	"database/sql"
	"fmt"
	"go-crud/internal/domain"

	"github.com/go-sql-driver/mysql"
)

func resolveSQLError(e error) error {
	if e == nil {
		return nil
	}
	if e == sql.ErrNoRows {
		return domain.ErrNotFound
	}
	if mysqlErr, ok := e.(*mysql.MySQLError); ok {
		switch mysqlErr.Number {
		case 1062:
			return domain.ErrAlreadyExists
		}
	}

	return fmt.Errorf("unexpected db error: %w", e)
}

