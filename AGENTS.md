# AGENTS.md - Go RESTful API (Clean Architecture)

## Commands

```bash
# Run server (requires MySQL running)
go run .

# Verify build (vet + format + build to temp, no repo artifacts)
go vet ./... && gofmt -l . && go build -o /tmp/api .

# Download deps (reconciles go.mod/go.sum — review changes before commit)
go mod tidy

# Generate Swagger docs after annotation changes
swag init -g main.go -o docs

# Run all tests
go test ./...

# Run a single test package (tests live in test/ subdirectory, NOT alongside packages)
go test ./test/service/... -v
go test ./test/repository/... -v
go test ./test/controller/... -v

# Run a single test by name
go test ./test/service -run TestName -count=1 -v
```

## CI/CD (`.github/workflows/go-ci.yml`)

- Runs on push to `main` and PRs to `main`.
- Order: `go mod download` → `go test ./...` → `go vet ./...` → `go build ./...` → `gofmt -l .`
- PR failures auto-comment via `actions/github-script`.
- `deploy` stage (needs `build`, only on `main`): deploys to **Railway** via `bervProject/railway-deploy@main` using `RAILWAY_TOKEN` secret.

## Architecture (strict layer boundaries)

| Layer | Package | Depends on | Must NOT contain |
|-------|---------|------------|------------------|
| HTTP | `controller/` | `service`, `payload`, `helper` | SQL, business logic |
| Business | `service/` | `repository`, `payload`, `entities`, `helper` | `net/http`, GORM queries |
| Persistence | `repository/` | `entities` | DTOs, HTTP codes, validation |
| DTOs | `payload/` | `entities` (mappers only) | DB/HTTP logic |
| Shared | `helper/`, `config/` | stdlib + drivers | Domain logic |

- **Test packages**: `test/controller/`, `test/repository/`, `test/service/` — use external test packages (`package X_test`) and are NOT co-located with source.

## Key Patterns

- **Controller**: parse input → call ONE service method → write ONE envelope. No method checks (Chi routes by method).
- **Service**: validate with `validator`, return `*helper.AppError`, map entities via `payload.New*Response`.
- **Repository**: GORM only, return `(*entities.*Entity, error)` / `(records, total, error)`; return package sentinel errors (e.g., `repository.ErrNotFound`) on not-found — use `errors.Is`, never compare with fresh `errors.New(...)`.
- **Never expose GORM models over HTTP** — always map to `payload` DTOs.

## Adding a New Resource

1. Create entity in `entities/`
2. Create payload (request/response) in `payload/`
3. Create repository in `repository/`
4. Create service in `service/`
5. Create controller in `controller/`
6. Register in `main.go`: wire repo → service → controller, add route under `/api/v1/`
7. Add `AutoMigrate(&entities.NewEntity{})` in `config.Migrate`

## Configuration

- `.env` is **ignored and optional** (both `.env` and `.http` are in `.gitignore`). Missing file uses defaults.
- Server port: `APP_PORT` (default `6767`)
- DB: `DB_USER`, `DB_PASSWORD`, `DB_HOST`, `DB_PORT`, `DB_NAME`
- Pool: `DB_MAX_IDLE_CONNS` (10), `DB_MAX_OPEN_CONNS` (100)
- `APP_URL` (optional): when set, overrides Swagger `Host`/`Schemes` at startup.

## Database

- **`config.Migrate` migrates ONLY `ProductEntity`** (config/config.go:78). Books and categories tables are NOT auto-created. You must add them manually or add `&entities.BookEntity{}` / `&entities.CategoryEntity{}` to `Migrate`.
- For production: use versioned migrations (`golang-migrate`).
- Connection: `charset=utf8mb4&parseTime=True&loc=Local`, 5m max-idle, 30m max-lifetime, ping on startup.

## API

- Versioned: `/api/v1/products`, `/api/v1/books`, `/api/v1/categories`
- **No legacy `/products` routes exist** — only versioned endpoints
- Pagination: `page` (default 1), `per_page` (default 10, max 100); categories are unpaginated
- Envelope: `{ "message": "...", "data": ..., "meta": {...} }`
- Errors: `{ "error": "..." }` with HTTP codes: 201, 400, 404, 500
- Swagger UI: `/swagger/index.html`

## Testing

- Tests are in `test/` subdirectory (NOT alongside source packages). Run via `go test ./test/...`.
- Repository tests use `sqlmock` (`github.com/DATA-DOG/go-sqlmock`) for DB mocking.
- Service/controller tests use `testify/mock` for dependency mocking.
- Test files: `package X_test` (external test packages).

## Validation Rules (service layer)

### Product (create/update)
| Field | Rule |
|-------|------|
| `name` | required, 1–255 chars |
| `description` | optional, max 5000 chars |
| `price` | required, `> 0` |
| `stock` | `>= 0` |

### Book (create/update)
| Field | Rule |
|-------|------|
| `title` | required |
| `category_id` | required |
| `author` | required |
| `stock` | required |

### Category (create/update)
- `CategoryRequest` has **no validation tags**; only JSON decoding is performed.

Validation returns first failure as `400` with message like `{"error":"Name is required"}`.

## Middleware & Lifecycle

- Chi middleware: `RequestID`, `RealIP`, `Logger`, `Recoverer`, `Timeout(60s)`
- Server timeouts: Read 15s, Write 15s, Idle 60s
- Graceful shutdown on `SIGINT`/`SIGTERM` with 10s deadline

## Current Gotchas

1. **`config.Migrate` migrates only `ProductEntity`** — books/categories tables not auto-created. Add entities to `Migrate` before expecting tables.
2. **`NewBookResponse` (payload/BookPayload.go:41) dereferences `BookEntity.CategoryId` (*uint) without nil check** — will panic if the entity has a nil `CategoryId`. Guard before calling or check at mapping time.
3. **`helper.ParseIDParam` always returns error message `"invalid product id"`** regardless of which resource/controller uses it (helper/response.go:50-57). Acceptable for products, misleading for books/categories.
4. **`ProductService.NewProductServiceImplWithDB`** is a deprecated backward-compat constructor (ignored first arg). Prefer `NewProductServiceImpl`.
