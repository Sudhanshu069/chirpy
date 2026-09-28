# Chirpy

A small Twitter-style REST API written in Go. Users can sign up and post short messages ("chirps") of up to 140 characters.

Built while working through the [Boot.dev](https://www.boot.dev) backend course.

## Stack

- **Go** standard library `net/http` (method and path-parameter routing, no framework)
- **PostgreSQL**, with **goose** for migrations and **sqlc** for type-safe generated queries
- **argon2id** password hashing (in progress)

## Run locally

Requires Go, PostgreSQL, and [goose](https://github.com/pressly/goose).

```bash
# .env
DB_URL="postgres://user:pass@localhost:5432/chirpy?sslmode=disable"
PLATFORM="dev"
```

```bash
goose -dir sql/schema postgres "$DB_URL" up
go run .
```

The server listens on `:8080`.
