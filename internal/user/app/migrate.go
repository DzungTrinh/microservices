//go:build migrate

package app

import (
	"errors"
	"fmt"
	"microservices/user-management/pkg/logger"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const (
	defaultAttempts   = 10
	defaultTimeout    = 2 * time.Second
	migrationFilePath = "db/migrations"
)

func init() {
	dsn, ok := os.LookupEnv("DATABASE_DSN")
	if !ok || dsn == "" {
		logger.GetInstance().Fatalf("migrate: DATABASE_DSN environment variable not set")
	}

	var (
		attempts = defaultAttempts
		err      error
		m        *migrate.Migrate
	)

	for attempts > 0 {
		inDocker, ok := os.LookupEnv("IN_DOCKER")
		if !ok || len(inDocker) == 0 {
			logger.GetInstance().Fatalf("migrate: environment variable not declared: IN_DOCKER")
		}

		dir := fmt.Sprintf("file://%s", migrationFilePath)

		if dockered, _ := strconv.ParseBool(inDocker); !dockered {
			cur, _ := os.Getwd()
			dir = fmt.Sprintf("file://%s/%s", filepath.Dir(cur+"/../../.."), migrationFilePath)
		}

		m, err = migrate.New(dir, dsn)
		if err == nil {
			break
		}

		logger.GetInstance().Printf("Migration: MySQL is trying to connect, attempts left: %d, error: %v", attempts, err)
		time.Sleep(defaultTimeout)
		attempts--
	}

	if err != nil {
		logger.GetInstance().Fatalf("Migration: MySQL connect error: %v", err)
	}

	defer m.Close()

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.GetInstance().Fatalf("Migration: up error: %v", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		logger.GetInstance().Printf("Migration: no change")
	} else {
		logger.GetInstance().Printf("Migration: up success")
	}
}
