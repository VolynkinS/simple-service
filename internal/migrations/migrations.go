package migrations

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Migration struct {
	Version     int
	Description string
	Query       string
}

var migrations = []Migration{
	{
		Version:     1,
		Description: "Create tasks table",
		Query: `
			CREATE TABLE IF NOT EXISTS tasks (
				id SERIAL PRIMARY KEY,
				title TEXT NOT NULL,
				description TEXT,
				status TEXT CHECK (status IN ('new', 'in_progress', 'done')) DEFAULT 'new',
				created_at TIMESTAMP DEFAULT now(),
				updated_at TIMESTAMP DEFAULT now()
			)`,
	},
	{
		Version:     2,
		Description: "Create schema_migrations table",
		Query: `
			CREATE TABLE IF NOT EXISTS schema_migrations (
				version INTEGER PRIMARY KEY,
				applied_at TIMESTAMP DEFAULT now()
			)`,
	},
}

func RunMigrations(ctx context.Context, pool *pgxpool.Pool, logger *zap.SugaredLogger) error {
	// Создаем таблицу миграций если её нет
	_, err := pool.Exec(ctx, migrations[1].Query)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	for _, migration := range migrations {
		if migration.Version == 2 {
			continue // Пропускаем создание таблицы миграций
		}

		// Проверяем, применена ли миграция
		var count int
		err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE version = $1", migration.Version).Scan(&count)
		if err != nil {
			return fmt.Errorf("failed to check migration status: %w", err)
		}

		if count > 0 {
			logger.Infof("Migration %d already applied, skipping", migration.Version)
			continue
		}

		// Применяем миграцию
		logger.Infof("Applying migration %d: %s", migration.Version, migration.Description)
		_, err = pool.Exec(ctx, migration.Query)
		if err != nil {
			return fmt.Errorf("failed to apply migration %d: %w", migration.Version, err)
		}

		// Записываем в таблицу миграций
		_, err = pool.Exec(ctx, "INSERT INTO schema_migrations (version) VALUES ($1)", migration.Version)
		if err != nil {
			return fmt.Errorf("failed to record migration %d: %w", migration.Version, err)
		}

		logger.Infof("Migration %d applied successfully", migration.Version)
	}

	return nil
}
