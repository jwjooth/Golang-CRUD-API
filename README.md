# Product CRUD API

Clean-architecture Product CRUD API built with Go, Chi, and GORM (MySQL).

Thin HTTP handlers → validating service layer → entity-only repository. Consistent JSON envelopes, versioned routes, pagination, and graceful shutdown.

## Features

- Full Product CRUD: list (paginated), get by ID, create, full update (PUT), delete
- Clean architecture with strict layer boundaries
- Request validation (`go-playground/validator/v10`) in the service layer
- Consistent success/error JSON envelopes
- Versioned API (`/api/v1`) with backward-compatible legacy `/products` routes
- Pagination metadata (`page`, `per_page`, `total`, `total_pages`)
- Chi middleware stack: `RequestID`, `RealIP`, `Logger`, `Recoverer`, `Timeout`
- MySQL connection pooling, health check, `AutoMigrate`, graceful shutdown
- Correct HTTP semantics: `201` on create, `400` on validation, `404` on missing rows

## Tech stack

| Component | Choice |
|---|---|
| Language | Go 1.25 |
| Router | `github.com/go-chi/chi/v5` |
| ORM | `gorm.io/gorm` + `gorm.io/driver/mysql` |
| Validation | `github.com/go-playground/validator/v10` |
| Env loading | `github.com/joho/godotenv` (optional `.env`) |
| Database | MySQL 8.x |

## Architecture

```
HTTP request
  → controller (decode params/body, encode JSON, no business logic)
  → service    (validation, business rules, entity ↔ DTO mapping)
  → repository (persistence only, entities in/out, no DTOs or HTTP codes)
  → MySQL (via GORM)
```

Layer rules:

| Layer | Package | May depend on | Must NOT contain |
|---|---|---|---|
| HTTP | `controller/` | `service`, `payload`, `helper` | SQL, business rules |
| Business | `service/` | `repository`, `payload`, `entities`, `helper` | `net/http`, GORM queries |
| Persistence | `repository/` | `entities` | DTOs, HTTP codes, validation messages |
| DTOs | `payload/` | `entities` (mappers only) | DB or HTTP logic |
| Shared | `helper/`, `config/` | stdlib + drivers | Domain logic |

## Project structure

```
.
├── main.go                     # wiring, middleware, routes, server lifecycle
├── config/
│   └── config.go               # env Config, OpenDB pool, Migrate, CloseDB
├── controller/
│   └── ProductController.go    # List, GetByID, Create, Update, Delete handlers
├── service/
│   └── ProductService.go       # validation + entity↔DTO mapping
├── repository/
│   └── ProductRepository.go    # Create/FindAll/FindByID/Update/Delete entities
├── entities/
│   └── ProductEntity.go        # GORM Product model, table `products`
├── payload/
│   └── ProductPayload.go       # Create/Update requests, ProductResponse, page meta
├── helper/
│   ├── error.go                # AppError{Code,Message,Err} + BadRequest/NotFound/Internal
│   └── response.go             # WriteSuccess/WriteSuccessWithMeta/WriteError envelopes
├── .env.example
└── go.mod / go.sum
```

## Requirements

- Go 1.25+
- MySQL 8.x running locally or remotely
- `curl` (for manual API checks)

## Quickstart

```bash
# 1. Clone and enter the repo
git clone <repo-url>
cd golang-restful-api

# 2. Configure environment
cp .env.example .env
# edit .env if needed

# 3. Create the database (once)
mysql -h localhost -P 3306 -u root -p -e "CREATE DATABASE IF NOT EXISTS products CHARACTER SET utf8mb4;"

# 4. Download modules
go mod tidy

# 5. Run
go run .
```

Server listens on `:6767` by default (`APP_PORT`).

Verify:

```bash
curl -s http://localhost:6767/healthz
# {"message":"ok","data":{"status":"up"}}

curl -s "http://localhost:6767/api/v1/products?page=1&per_page=5"
```

Build a binary:

```bash
go vet ./...
go build -o api .
./api
```

## Configuration

`.env` is optional. Missing file falls back to env vars/defaults (production-safe).

| Variable | Default | Description |
|---|---|---|
| `APP_PORT` | `6767` | HTTP listen port |
| `DB_USER` | `root` | MySQL user |
| `DB_PASSWORD` | `` | MySQL password |
| `DB_HOST` | `localhost` | MySQL host |
| `DB_PORT` | `3306` | MySQL port |
| `DB_NAME` | `products` | MySQL database |
| `DB_MAX_IDLE_CONNS` | `10` | Connection pool idle limit |
| `DB_MAX_OPEN_CONNS` | `100` | Connection pool open limit |

Connection uses `charset=utf8mb4&parseTime=True&loc=Local`, 5 min max-idle time, 30 min max lifetime, plus a ping check on startup.

## Database

Schema is created via `config.Migrate` (`AutoMigrate`) on boot. Suitable for development; use versioned migrations (e.g. `golang-migrate`) for production evolution.

`products` table (effective schema):

```sql
CREATE TABLE products (
  id          INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  name        VARCHAR(255) NOT NULL,
  description TEXT NULL,
  price       DOUBLE NOT NULL,
  stock       INT NOT NULL DEFAULT 0,
  created_at  DATETIME NULL,
  updated_at  DATETIME NULL
);
```

Entity (`entities/ProductEntity.go`):

```go
type ProductEntity struct {
  ID          uint      `json:"id"`
  Name        string    `json:"name"`
  Description string    `json:"description"`
  Price       float64   `json:"price"`
  Stock       int       `json:"stock"`
  CreatedAt   time.Time `json:"created_at"`
  UpdatedAt   time.Time `json:"updated_at"`
}
```

## API reference

Base URL: `http://localhost:6767`

- Versioned (current): `/api/v1/products`
- Legacy (deprecated, same behavior): `/products`
- All bodies are JSON. All responses set `Content-Type: application/json`.

| Method | Path | Description | Success |
|---|---|---|---|
| `GET` | `/healthz` | Liveness probe | `200` |
| `GET` | `/api/v1/products?page=&per_page=` | List products (paginated) | `200` |
| `GET` | `/api/v1/products/{id}` | Get one product | `200` |
| `POST` | `/api/v1/products` | Create product | `201` |
| `PUT` | `/api/v1/products/{id}` | Full product replacement | `200` |
| `DELETE` | `/api/v1/products/{id}` | Delete product | `200` |

Pagination: `page` defaults to `1`, `per_page` defaults to `10`, capped at `100`. Invalid values fall back to defaults.

### Envelopes

Success:

```json
{ "message": "success", "data": {} }
```

Paginated list:

```json
{
  "message": "success",
  "data": [ { "id": 1, "name": "...", "price": 129.99, "stock": 48, "...": "..." } ],
  "meta": { "page": 1, "per_page": 5, "total": 10, "total_pages": 2 }
}
```

Delete success:

```json
{ "message": "product deleted successfully" }
```

Error (all failures):

```json
{ "error": "product not found" }
```

| HTTP | Meaning | Example |
|---|---|---|
| `201` | Created | POST success |
| `400` | Bad request | invalid JSON, invalid `id`, validation failure |
| `404` | Not found | unknown product ID on get/update/delete |
| `500` | Internal | DB failure |

### Product object

```json
{
  "id": 12,
  "name": "Test Keyboard",
  "description": "smoke test",
  "price": 99.99,
  "stock": 10,
  "created_at": "2026-09-11T15:58:04.028+07:00",
  "updated_at": "2026-09-11T15:58:04.028+07:00"
}
```

### Validation rules

Applied in the service layer for both create and PUT:

| Field | Rule |
|---|---|
| `name` | required, 1–255 chars |
| `description` | optional, max 5000 chars |
| `price` | required, `> 0` |
| `stock` | `>= 0` (zero allowed = out of stock) |

Validation failures return `400` with the first failing reason, e.g. `{"error":"Name is required"}`.

### curl examples

```bash
BASE=http://localhost:6767/api/v1

# List (page 1, 5 per page)
curl -s "$BASE/products?page=1&per_page=5"

# Get one
curl -s "$BASE/products/1"

# Create
curl -s -X POST "$BASE/products" \
  -H 'Content-Type: application/json' \
  -d '{"name":"USB-C GaN Charger 65W","description":"65W fast charger","price":49.9,"stock":125}'

# Full update (PUT requires all fields)
curl -s -X PUT "$BASE/products/12" \
  -H 'Content-Type: application/json' \
  -d '{"name":"Test Keyboard v2","description":"updated","price":119.5,"stock":5}'

# Delete
curl -s -X DELETE "$BASE/products/12"

# Error cases
curl -s "$BASE/products/abc"        # {"error":"invalid product id"}
curl -s "$BASE/products/999999"     # {"error":"product not found"}
```

## Middleware and lifecycle

- `RequestID`, `RealIP`, `Logger`, `Recoverer`, `Timeout(60s)` on every request
- `ReadTimeout: 15s`, `WriteTimeout: 15s`, `IdleTimeout: 60s`
- `SIGINT`/`SIGTERM` trigger a graceful shutdown with a 10s deadline

## Development

```bash
go mod tidy
go vet ./...
gofmt -l .
go build -o /tmp/api-build .
```

Conventions:

- Controller: parse input, call one service method, write one envelope. No `if r.Method != ...` checks (Chi routes by method).
- Service: validate with `validator`, return `*helper.AppError`, map entities via `payload.NewProductResponse`.
- Repository: GORM only, return `(*entities.ProductEntity, error)` / `(records, total, error)`; return `repository.ErrNotFound` when `RowsAffected == 0` or `gorm.ErrRecordNotFound`.
- Never expose GORM models directly over HTTP; always map to `payload` DTOs.

Adding a new resource (e.g. `orders`): copy the `entities → payload → repository → service → controller` chain, register `/api/v1/orders` in `main.go`, add `AutoMigrate(&entities.OrderEntity{})`.

## Deployment notes

- Do not require `.env` in production; inject real env vars instead.
- Put the service behind a reverse proxy / TLS terminator for public traffic.
- Replace `AutoMigrate` with versioned migrations if the schema is shared or long-lived.
- Set `DB_MAX_OPEN_CONNS` according to MySQL `max_connections` and instance count.

## License

No license file is currently included. Add one (e.g. MIT) before public distribution.
