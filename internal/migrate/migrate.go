package migrate

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
)

func ApplyMigrations(db *sql.DB) error {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return fmt.Errorf("failed to create MySQL driver instance: %w", err)
	}
	defer driver.Close()

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"mysql",
		driver,
	)
	fmt.Println(m)
	if err != nil {
		return fmt.Errorf("migration init error: %w", err)
	}

	err = m.Up()
	switch err {
	case nil:
		log.Println("migrations applied successfully")
	case migrate.ErrNoChange:
		log.Println("no migrations to apply")
	default:
		return fmt.Errorf("migration up error: %w", err)
	}
	return nil
}
