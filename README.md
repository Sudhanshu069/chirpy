# Chirpy

A Twitter-style REST API in Go. Users sign up, log in, and post short messages ("chirps") of up to 140 characters.

Built with nothing but Go's standard library `net/http`. No web framework, so every piece of auth, routing, and error handling is visible and easy to follow.

## Why it's interesting

- **Real auth flow:** argon2id password hashing, short-lived JWT access tokens, and long-lived refresh tokens that can be revoked.
- **Ownership checks:** users can only edit their own account and delete their own chirps (`403` otherwise).
- **Webhooks:** a payment-provider webhook upgrades users to "Chirpy Red", authenticated with an API key.
- **Type-safe SQL:** queries are plain SQL compiled to Go with [sqlc](https://sqlc.dev); schema changes are versioned with [goose](https://github.com/pressly/goose).

## Run it

Requires Go 1.26+, PostgreSQL, and goose.

1. Create a `.env` file in the project root:

   ```bash
   DB_URL="postgres://user:pass@localhost:5432/chirpy?sslmode=disable"
   PLATFORM="dev"            # enables POST /admin/reset
   JWT_SECRET="..."          # generate with: openssl rand -hex 32
   POLKA_KEY="..."           # API key expected from the webhook sender
   ```

2. Run the migrations and start the server:

   ```bash
   set -a && . ./.env && set +a
   goose -dir sql/schema postgres "$DB_URL" up
   go run .
   ```

The server listens on `http://localhost:8080`.

## API

| Method | Path | Auth | Description |
|---|---|---|---|
| `POST` | `/api/users` | | Create a user |
| `PUT` | `/api/users` | JWT | Update your email and password |
| `POST` | `/api/login` | | Get an access token and a refresh token |
| `POST` | `/api/refresh` | Refresh token | Get a new access token |
| `POST` | `/api/revoke` | Refresh token | Revoke a refresh token |
| `POST` | `/api/chirps` | JWT | Post a chirp |
| `GET` | `/api/chirps` | | List chirps. Optional `?author_id=<uuid>` and `?sort=asc\|desc` |
| `GET` | `/api/chirps/{chirpID}` | | Get one chirp |
| `DELETE` | `/api/chirps/{chirpID}` | JWT | Delete your own chirp |
| `POST` | `/api/polka/webhooks` | API key | Upgrade a user to Chirpy Red |
| `GET` | `/api/healthz` | | Health check |

Authenticated requests send `Authorization: Bearer <token>`, or `Authorization: ApiKey <key>` for the webhook.

---

Built while working through the [Boot.dev](https://www.boot.dev) backend course.
