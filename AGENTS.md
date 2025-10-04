# Repository Guidelines

## Project Structure & Module Organization
Simple Service is a Go REST API built on Fiber with PostgreSQL persistence. Core layout:
- `cmd/main.go` boots the HTTP server, loads env config, and registers routes.
- `internal/api` holds handlers and middleware; align new feature folders with existing modules.
- `internal/service` encapsulates business rules and relies on `internal/repo` for pgx-based data access.
- `pkg/validator` hosts reusable helpers with dedicated tests in `pkg/validator/validator_test.go`.
- `migrations/postgres` stores ordered SQL migrations applied automatically by Docker.
- `docs` captures generated Swagger artifacts; regenerate instead of editing manually.

## Build, Test, and Development Commands
- `docker-compose up -d postgres` starts PostgreSQL with seed migrations.
- `make deps` syncs Go modules.
- `make build` produces `bin/simple-service`.
- `make run` launches the API using `local.env`.
- `make test` runs `go test -v ./... -cover`.
- `make lint` executes `golangci-lint run`; ensure the binary is on PATH.
- `make swagger-gen` refreshes OpenAPI docs from `cmd/main.go` annotations.

## Coding Style & Naming Conventions
- Format with `go fmt ./...`; do not commit unformatted code.
- Keep lint clean; follow the guidance emitted by `golangci-lint`.
- Use UpperCamelCase for exported identifiers and lowerCamelCase for internals; file names stay snake_case.
- Group feature-specific DTOs under `internal/dto` and keep middleware small and focused.

## Testing Guidelines
- Prefer table-driven tests in `_test.go` files alongside the code they verify.
- Maintain or improve coverage reported by `make test`; add service and repo assertions when adjusting business rules.
- Seed integration tests against Docker DB via `local.env` credentials; clean up created rows.

## Commit & Pull Request Guidelines
- Follow Conventional Commits (`feat:`, `fix:`, `chore:`) as seen in `git log`.
- Keep commits scoped narrowly and include related migrations or docs updates.
- PRs need a short summary, testing notes (e.g., `make test`), and linked issues or TODO references; attach Swagger screenshots when endpoints change.
- Request review only after lint and tests pass locally.

## Environment & Configuration Tips
- Copy `local.env` as the base for new environment files; never commit secrets.
- Update `Dockerfile` and `docker-compose.yml` together when adding config variables.
- Document new env vars in `README.md` and ensure sensible defaults exist for local development.
