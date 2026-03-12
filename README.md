## Chirpy Backend

Chirpy is a **Twitter‑like mini backend** implemented in Go. It exposes a JSON HTTP API for posting short messages (“chirps”), managing users, and handling authentication with JWT access tokens and refresh tokens.

This project was built as part of the **Boot.dev Chirpy course**, and serves in this portfolio as a demonstration of backend fundamentals: HTTP APIs, authentication, authorization, relational data modeling, and working with PostgreSQL from Go.

## Features

- **User management**
  - Sign‑up with email and password
  - Login with secure password verification
  - Update email and password (authenticated)
  - `is_chirpy_red` flag for premium users

- **Authentication & authorization**
  - Password hashing using **Argon2id**
  - **JWT** access tokens with HMAC signing
  - Long‑lived, revocable **refresh tokens**
  - Bearer token protection on mutating endpoints

- **Chirps (short messages)**
  - Create, list, retrieve, and delete chirps
  - 140‑character limit with server‑side validation
  - Simple profanity filtering on specific words
  - Filtering by author and optional descending sort

- **Admin & observability**
  - Health/readiness endpoint
  - Simple HTML metrics page that tracks static file hits
  - Development‑only reset endpoint to clear users and metrics

- **Webhook integration**
  - “Polka” webhook endpoint to **upgrade users to Chirpy Red** based on signed events and API keys

## Tech Stack

- **Language**: Go (tested with **Go 1.21+**)
- **Framework**: Standard library `net/http`
- **Database**: PostgreSQL
- **Database access**: 
  - Raw SQL + **sqlc**‑generated type‑safe Go code (`internal/database`)
  - SQL schema files in `sql/schema` and query definitions in `sql/queries`
- **Auth & security libraries**:
  - `github.com/alexedwards/argon2id` for password hashing
  - `github.com/golang-jwt/jwt/v5` for JWT creation and validation
- **Migrations**: SQL migrations managed with **Goose** (applied from `sql/schema`)

## High‑Level Architecture

- `main.go` wires together:
  - Environment configuration (database URL, secrets, platform)
  - Database connection and `sqlc` query struct
  - HTTP routes and handlers on `net/http.ServeMux`
- `internal/database` contains the `sqlc`‑generated models and query methods.
- `internal/auth` encapsulates:
  - Password hashing and verification (`hash.go`)
  - JWT creation and validation (`jwt.go`)
  - Bearer/API key extraction helpers and refresh token generation.
- Top‑level handler files (`users.go`, `chirps.go`, `refresh.go`, etc.) implement the JSON API:
  - Decode/validate requests
  - Call into the database layer
  - Return typed JSON responses and consistent error shapes

This structure keeps HTTP concerns, auth logic, and persistence concerns clearly separated while still remaining small and easy to navigate.

## Getting Started (High‑Level)

This project is primarily meant to demonstrate backend skills. The instructions below are intentionally high‑level; adapt them as needed for your environment.

### Prerequisites

- Go **1.21+**
- PostgreSQL (any recent version such as 14+)
- Access to a tool to run SQL migrations (e.g. **Goose**)

### Environment

The server expects the following environment variables:

- **`DB_URL`** – PostgreSQL connection string  
  Example: `postgres://user:password@localhost:5432/chirpy?sslmode=disable`
- **`SECRET_KEY`** – HMAC secret used to sign JWT access tokens
- **`POLKA_KEY`** – shared API key used to authenticate incoming webhook requests
- **`PLATFORM`** – environment name (e.g. `dev`); controls access to certain admin endpoints like `/admin/reset`

You can manage these via a local `.env` file (the project uses `github.com/joho/godotenv` in `main.go`).

### Database & Migrations

1. **Create a database** in PostgreSQL, for example:
   - `CREATE DATABASE chirpy;`
2. **Run the SQL migrations** in `sql/schema` in order using Goose or your preferred migration tool.  
   These migrations create:
   - `users` (including `is_chirpy_red` and `hashed_password`)
   - `chirps`
   - `refresh_tokens`
3. Confirm that `DB_URL` points at this database.

### Running the Server

With the environment configured and database migrated:

- Run the server from the repository root:

```bash
go run ./...
```

The HTTP server listens on **port 8080** by default and serves:

- Static content from `index.html` at `/app`
- JSON API under `/api/*`

## API Overview

The API is designed to be simple to consume from any frontend or API client.

### Users & Authentication

- **`POST /api/users`** – create a new user  
  - Body: email and password  
  - Response: user object (without password) with `id`, timestamps, and `is_chirpy_red`

- **`POST /api/login`** – login and obtain tokens  
  - Body: email and password  
  - Response: 
    - User details
    - Short‑lived JWT **access token** (`token`)
    - Long‑lived **refresh token** (`refresh_token`)

- **`PUT /api/users`** – update current user’s email and password  
  - Requires `Authorization: Bearer <access-token>`  
  - Body: new email and password  
  - Response: updated user object

- **`POST /api/refresh`** – exchange a refresh token for a new access token  
  - Requires `Authorization: Bearer <refresh-token>`  
  - Response: new access token

- **`POST /api/revoke`** – revoke a refresh token  
  - Requires `Authorization: Bearer <refresh-token>`  
  - Marks the refresh token as revoked and returns `204 No Content`

### Chirps

- **`POST /api/chirps`** – create a chirp  
  - Requires `Authorization: Bearer <access-token>`  
  - Server enforces a **140‑character limit** and filters a small list of forbidden words.

- **`GET /api/chirps`** – list chirps  
  - Optional query params:
    - `author_id` – filter by user ID
    - `sort=desc` – return in reverse chronological order

- **`GET /api/chirps/{chirpID}`** – get a single chirp by ID

- **`DELETE /api/chirps/{chirpID}`** – delete a chirp  
  - Requires `Authorization: Bearer <access-token>`  
  - Only the author of the chirp can delete it.

### Webhooks & Admin

- **`POST /api/polka/webhooks`** – handle “Polka” webhook events  
  - Secured with `Authorization: ApiKey <POLKA_KEY>`  
  - On `user.upgraded` events, upgrades the specified user to **Chirpy Red**.

- **`GET /api/healthz`** – simple readiness/health endpoint.

- **`GET /admin/metrics`** – basic HTML metrics page showing file‑server hit count.

- **`POST /admin/reset`** – development‑only reset of metrics and users  
  - Only active when `PLATFORM` is set to `dev`.

## Skills Demonstrated

- **HTTP API design in Go**
  - Routing with `net/http.ServeMux` and clear separation of handlers.
  - JSON request/response handling and structured error responses.

- **Authentication & security**
  - Argon2id password hashing and verification.
  - JWT‑based access tokens with HMAC signing.
  - Refresh token model with expiry and explicit revocation.
  - Bearer and API key authorization schemes for different endpoints.

- **Relational data modeling & SQL**
  - Normalized tables for users, chirps, and refresh tokens.
  - Use of `sqlc` for type‑safe query generation.
  - SQL migrations and schema evolution with Goose.

- **Backend architecture & maintainability**
  - Clear separation between HTTP layer, auth utilities, and database layer.
  - Use of configuration via environment variables for secrets and runtime behavior.
  - Small but meaningful admin and observability endpoints (health, metrics, reset).

## License

This project is provided as a **portfolio / learning project**. You are welcome to read and learn from the code; adapt or reuse parts of it at your own risk and with appropriate attribution.

