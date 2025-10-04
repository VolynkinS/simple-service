# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Architecture

Simple Service is a Go REST API using Fiber framework with PostgreSQL persistence, organized in clean architecture layers:

**Layer Flow**: HTTP Request → Handler → Service (business logic) → Repository (data access) → PostgreSQL

- **cmd/main.go**: Application entry point that loads `local.env`, initializes components, runs migrations, and starts the Fiber server with graceful shutdown
- **internal/api**: HTTP layer with `api.go` (route registration), `handlers/` (request/response handling), `middleware/` (JWT auth via Bearer token)
- **internal/service**: Business logic layer defining `Service` interface and `Repository` interface (dependency inversion)
- **internal/repo**: Data access layer implementing repository with pgxpool connection management
- **internal/dto**: HTTP request/response structs with validation tags and Swagger annotations
- **internal/config**: Configuration structs loaded via envconfig from environment variables
- **internal/logger**: Zap logger wrapper initialization
- **internal/migrations**: Embedded SQL migrations run automatically at startup
- **pkg/validator**: Reusable validation utilities (e.g., title/description length checks)
- **docs/**: Auto-generated Swagger documentation artifacts (never edit manually)

**Key patterns**:
- Dependency injection via constructor functions (`NewService`, `NewRepository`, `NewTaskHandler`)
- Interface-based dependencies: service depends on `Repository` interface, not concrete `repository` struct
- Context propagation through all layers for cancellation and timeouts
- Structured logging with zap.SugaredLogger

## Development Commands

### Environment Setup
```bash
# Copy and configure local environment (DB credentials must be filled in local.env)
cp local.env .env  # if needed
# Edit local.env to set DB_HOST, DB_PORT, DB_NAME, DB_USER, DB_PASSWORD
```

### Database
```bash
# Start PostgreSQL with migrations via docker-compose
docker-compose up -d postgres

# Full stack (app + postgres)
docker-compose up -d

# Stop and remove containers
docker-compose down
```

### Build & Run
```bash
# Install dependencies
make deps

# Build binary to bin/simple-service
make build

# Run locally (reads local.env, requires postgres running)
make run

# Run in IDE debugger (recommended for development)
# Use GoLand/VSCode run configurations pointing to cmd/main.go
```

### Testing
```bash
# Run all tests with coverage
make test

# Run specific package tests
go test -v ./internal/service/... -cover

# Run single test
go test -v ./pkg/validator -run TestValidateTitle
```

### Code Quality
```bash
# Format code
make fmt

# Run linter (requires golangci-lint installed)
make lint

# Full build pipeline
make all  # deps + swagger-gen + build
```

### Swagger Documentation
```bash
# Install swag CLI (one-time)
make swagger-install

# Regenerate docs from annotations
make swagger-gen

# View Swagger UI (after starting server)
# Navigate to http://localhost:8080/swagger/
```

## Configuration

All configuration is loaded from environment variables via envconfig. The `local.env` file provides defaults for local development.

**Required variables**:
- `DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASSWORD`: PostgreSQL connection
- `PORT`: HTTP server listen address (e.g., `:8080`)
- `TOKEN`: Bearer token for JWT middleware authorization
- `LOG_LEVEL`: Logging level (`debug`, `info`, `warn`, `error`)
- `SERVER_NAME`: Server identifier for logging
- `WRITE_TIMEOUT`: HTTP write timeout duration

**Optional with defaults**:
- `DB_SSL_MODE`: Default `disable`
- `DB_POOL_MAX_CONNS`: Default `5`
- `DB_POOL_MAX_CONN_LIFETIME`: Default `180s`
- `DB_POOL_MAX_CONN_IDLE_TIME`: Default `100s`

See [internal/config/config.go](internal/config/config.go) for complete config structure.

## Adding New Endpoints

1. **Define DTOs** in `internal/dto/` with validation tags and Swagger annotations:
```go
// CreateTaskRequest description
// @Description Request body for creating a task
type CreateTaskRequest struct {
    Title       string `json:"title" validate:"required" example:"Implement feature"`
    Description string `json:"description" example:"Add new API endpoint"`
} // @name CreateTaskRequest
```

2. **Add Service method** in `internal/service/service.go`:
```go
// Update Service interface
type Service interface {
    CreateTask(ctx context.Context, req TaskRequest) (int, error)
    NewMethod(ctx context.Context, params SomeType) (*Result, error)
}

// Implement method on service struct
func (s *service) NewMethod(ctx context.Context, params SomeType) (*Result, error) {
    // Business logic here
    return s.repo.SomeRepoMethod(ctx, params)
}
```

3. **Add Repository method** if needed in `internal/repo/repo.go`:
```go
// Define SQL query as const
const someQuery = `SELECT ... FROM ...`

// Implement method
func (r *repository) SomeRepoMethod(ctx context.Context, params SomeType) (*Result, error) {
    // pgxpool query execution
}
```

4. **Create Handler** in `internal/api/handlers/handlers.go`:
```go
// @Summary Brief description
// @Description Detailed description
// @Tags tasks
// @Accept json
// @Produce json
// @Param request body dto.RequestType true "Request description"
// @Success 200 {object} dto.ResponseType
// @Failure 400 {object} dto.ErrorResponse
// @Router /v1/endpoint [post]
func (h *TaskHandler) NewHandler(ctx *fiber.Ctx) error {
    // Parse request, validate, call service, return response
}
```

5. **Register route** in `internal/api/api.go`:
```go
apiGroup.Post("/endpoint", taskHandler.NewHandler)
```

6. **Regenerate Swagger**:
```bash
make swagger-gen
```

## Testing Guidelines

- Place tests alongside implementation (`service_test.go` next to `service.go`)
- Use table-driven tests for multiple cases
- Mock repository interface using `internal/repo/mocks/Repository.go` (generated via mockery)
- Integration tests should connect to Docker postgres via `local.env` credentials
- Clean up test data after integration tests
- Maintain or improve coverage shown in `make test` output

## Migrations

SQL migrations are in `migrations/postgres/` with sequential numbering:
- `000001_task.up.sql`: Create tasks table
- `000001_task.down.sql`: Drop tasks table

Migrations run automatically at startup via `internal/migrations/migrations.go`. For new migrations:
1. Create `000002_<name>.up.sql` and `000002_<name>.down.sql`
2. Add migration logic to `migrations.go` if needed
3. Restart app or run `docker-compose up -d postgres` to apply

## Authentication

All `/v1/*` routes require Bearer token authentication via `middleware.JWTAuthorization`.

Request header: `Authorization: Bearer <TOKEN_from_env>`

Swagger UI endpoint `/swagger/*` is public (no auth required).
