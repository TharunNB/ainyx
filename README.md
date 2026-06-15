# Ainyx User API

RESTful API to manage users with name and date of birth. Age is calculated dynamically on fetch.

**Stack:** Go · GoFiber · PostgreSQL · SQLC · Uber Zap · go-playground/validator

---

## Setup & Run

### Local

```bash
# 1. Copy env and fill in your DB credentials
cp .env.example .env

# 2. Install dependencies
go mod download

# 3. Run (migrations apply automatically on startup)
go run ./cmd/server/main.go
```

### Docker

```bash
docker compose up --build
```

Server starts on `http://localhost:8080`.

---

## Environment Variables

| Variable      | Default       | Description           |
|---------------|---------------|-----------------------|
| `APP_PORT`    | `8080`        | HTTP listen port      |
| `APP_ENV`     | `development` | `development` / `production` |
| `DB_HOST`     | `localhost`   | Postgres host         |
| `DB_PORT`     | `5432`        | Postgres port         |
| `DB_USER`     | `postgres`    | Postgres user         |
| `DB_PASSWORD` | —             | Postgres password     |
| `DB_NAME`     | `ainyx_db`    | Database name         |
| `DB_SSLMODE`  | `disable`     | SSL mode              |

---

## API Reference

All dates use ISO 8601 format: `YYYY-MM-DD`

### Health

| Method | Path      | Response        |
|--------|-----------|-----------------|
| GET    | `/health` | `{"status":"ok"}` |

---

### POST `/users` — Create user

**Request**
```json
{ "name": "Alice", "dob": "1990-05-10" }
```
**Response** `201 Created`
```json
{ "id": 1, "name": "Alice", "dob": "1990-05-10" }
```

---

### GET `/users/:id` — Get user by ID

**Response** `200 OK`
```json
{ "id": 1, "name": "Alice", "dob": "1990-05-10", "age": 35 }
```

---

### PUT `/users/:id` — Update user

**Request**
```json
{ "name": "Alice Updated", "dob": "1991-03-15" }
```
**Response** `200 OK`
```json
{ "id": 1, "name": "Alice Updated", "dob": "1991-03-15" }
```

---

### DELETE `/users/:id` — Delete user

**Response** `204 No Content`

---

### GET `/users` — List users (paginated)

**Query params:** `page` (default `1`), `limit` (default `10`, max `100`)

**Response** `200 OK`
```json
{
  "data": [{ "id": 1, "name": "Alice", "dob": "1990-05-10", "age": 35 }],
  "page": 1,
  "limit": 10,
  "total_items": 1,
  "total_pages": 1
}
```

---

## Error Response Shape

```json
{ "error": "<description>" }
```

| Status | Meaning                        |
|--------|--------------------------------|
| `400`  | Invalid path param or body     |
| `404`  | User not found                 |
| `422`  | Validation failure             |
| `500`  | Internal server error          |

---

## Tests

```bash
go test ./...
```
