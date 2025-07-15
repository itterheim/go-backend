# Exploration of Go Backend project structure

An exploration of

- Go project structure
- authentication with JWT
- using PostgreSQL
- SQL migrations

## Dependencies

- Go 1.24.24
- PGX v5 - github.com/jackc/pgx/v5
- uuid - github.com/google/uuid
- JWT - github.com/golang-jwt/jwt
- YAML / Viper - github.com/spf13/viper

## Folder structure

- `cmd/api` - API server
- `cmd/migrations` - migrations
- `config` - YAML configuration files
- `internal`
- `migrations` - SQL files for migrations
- `pkg/auth` - generic types for authorization, mostly for use in `/pkg`
- `pkg/cli` - CLI utilities
- `pkg/handler` - generic types and utilities for http handlers
- `pkg/middleware` - useful http middlewares
