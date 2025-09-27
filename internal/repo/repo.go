package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"

	"simple-service/internal/config"
	"simple-service/internal/service"
)

// Слой репозитория, здесь должны быть все методы, связанные с базой данных

// SQL-запросы
const (
	insertTaskQuery = `INSERT INTO tasks (title, description) VALUES ($1, $2) RETURNING id;`
	getTaskQuery    = `SELECT id, title, description, status, created_at, updated_at FROM tasks WHERE id = $1;`
)

type repository struct {
	pool *pgxpool.Pool
}

// NewRepository - создание нового экземпляра репозитория с подключением к PostgreSQL
func NewRepository(ctx context.Context, cfg config.PostgreSQL) (*repository, error) {
	// Формируем строку подключения
	connString := fmt.Sprintf(
		`user=%s password=%s host=%s port=%d dbname=%s sslmode=%s
        pool_max_conns=%d pool_max_conn_lifetime=%s pool_max_conn_idle_time=%s`,
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
		cfg.PoolMaxConns,
		cfg.PoolMaxConnLifetime.String(),
		cfg.PoolMaxConnIdleTime.String(),
	)

	// Парсим конфигурацию подключения
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse PostgreSQL config")
	}

	// Оптимизация выполнения запросов (кеширование запросов)
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe

	// Создаём пул соединений с базой данных
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create PostgreSQL connection pool")
	}

	return &repository{pool}, nil
}

// Close - закрытие пула соединений
func (r *repository) Close() {
	r.pool.Close()
}

// Pool - получение пула соединений
func (r *repository) Pool() *pgxpool.Pool {
	return r.pool
}

// CreateTask - вставка новой задачи в таблицу tasks
func (r *repository) CreateTask(ctx context.Context, task service.Task) (int, error) {
	var id int
	err := r.pool.QueryRow(ctx, insertTaskQuery, task.Title, task.Description).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "failed to insert task")
	}
	return id, nil
}

// GetTask - получение задачи по ID
func (r *repository) GetTask(ctx context.Context, id int) (*service.TaskResponse, error) {
	var task service.TaskResponse
	err := r.pool.QueryRow(ctx, getTaskQuery, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New("task not found")
		}
		return nil, errors.Wrap(err, "failed to get task")
	}
	return &task, nil
}
