package mysql

import (
	"database/sql"
	"microservices/user-management/pkg/logger"
	"time"
)

// NewDB initializes a MySQL database connection
func NewDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	logger.GetInstance().Printf("Database connection pool initialized")
	return db, nil
}
