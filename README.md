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

## Database

- `/migrations/00000000-00-functions.sql`
    - most of the tables should have an `created` and `updated` colums
    - currently just one resable function that handle updates of the `updated` column
- `/migrations/00000000-01-auth.sql`
    - there are two entities which can access the Rest API `users` and `devices`
    - `users`
        - they use JWT auth and refresh tokens
        - http only cookies
        - standard login with username and password
    - `tokens`
        - registry of user refresh tokens (its JTI and expiration)
        - for blacklisting
    - `devices`
        - custom ESP32 devices and other services reading and writing data to this API
        - uses a long lived JWT token in `authorization Bearer ...`
        - user has to create a device, generate a token
        - TODO: scope/role/permissions for a specific part of API (environment data, image rendering for eInk display, ...)
        - TODO: might potentially add support for access and refresh tokens
            1. device will register its identification
            2. user will approve the device
            3. long lived refresh token is generated
            4. device has one chance to retrieve its
            5. then it is the standard access/refresh procedure
