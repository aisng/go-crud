package server

import (
	"database/sql"
	"net/http"
)

func RunServer(db *sql.DB, handler http.Handler) error {
	if err := db.Ping(); err != nil {
		return err
	}

	server := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	return server.ListenAndServe()
}
