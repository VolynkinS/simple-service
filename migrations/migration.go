package migrations

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/pkg/errors"
	"simple-service/internal/config"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// ApplyMigrations - активирует неприменённые миграции из каталога /migrations/postgres
func ApplyMigrations(cfg config.PostgreSQL) error {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
	)

	m, err := migrate.New(
		"file://./migrations/postgres",
		dbURL,
	)
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return errors.Wrap(err, "failed to init migrations")
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return errors.Wrap(err, "failed to apply migrations")
	}

	return nil
}
