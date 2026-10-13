# AGENTS.md

Go + Chi + GORM (MySQL) Product REST API. Module `golang-restful-api`, Go 1.26. Entry: `main.go` (wiring, routes, graceful shutdown). No Dockerfile / compose in repo.

## Structure

Product-only clean architecture — `controller/` → `service/` → `repository/` → MySQL (`products` table):

- `controller/ProductController.go` — HTTP only: parse input via `helper.ParseIDParam` / `helper.QueryInt`, call service, write `helper.WriteSuccess*` / `WriteAppError`. No SQL, no business rules.
- `service/ProductService.go` — validation (go-playground/validator), business rules, entity↔DTO mapping. No `net/http`. Constructor takes repository interface.
- `repository/ProductRepository.go` — GORM only. Returns entities + `helper` error sentinels. No DTOs, no HTTP codes.
- `payload/ProductPayload.go` — request/response DTOs + mappers. `Product.Category` is a plain string column (≤50 chars), not a relation. `config/` — env load, DB pool, `Migrate` (ProductEntity only). `entities/ProductEntity.go` — GORM model. `helper/` — envelopes (`{"message","data"(,"meta")}` / `{"error":...}`), `AppError` → status mapping.
- Routes: `/healthz`, `/swagger/*`, `/api/v1/products` only (GET list paginated, POST, GET/PUT/DELETE `/{id}`). List is paginated (`?page=&per_page=`, `helper.QueryInt` falls back silently on bad input).

## Commands

```bash
cp .env.example .env   # .env optional; config.Load() falls back to defaults (APP_PORT=6767, DB localhost:3306/products)
go run .               # needs live MySQL; curl -i localhost:6767/healthz
go test ./...          # no DB needed — repo tests use go-sqlmock, service/controller use mocks
go test ./test/service/ -run TestXxx -v   # focused run
swag init -g main.go -o docs   # after changing swagger annotations/payloads; commit docs/
```

CI (`go-ci.yml`, runs on PR/push to `main`): `go test ./...` → `go vet ./...` → `go build ./...` → `gofmt -l` must be empty. Push to `main` auto-deploys to Railway (`Product-CRUD-API` service). Lint job runs golangci-lint (`errcheck, govet, staticcheck, unused` per `.golangci.yml`).

## Gotchas

- `config.Migrate` auto-migrates ProductEntity on boot — fine for dev, use versioned migrations in production.
- `docs/` is stale: still contains Book/Category/ProductCategory swagger specs (no Go toolchain here to run `swag init`). Regenerate with `swag init -g main.go -o docs` before trusting Swagger UI.
- `godotenv.Load()` failure is ignored; never rely on `.env` existing in prod — set env directly. Never commit `.env`.
- `APP_URL` (e.g. `https://api.example.com`) sets Swagger host/scheme; empty default serves localhost docs.
- Server has 60s chi timeout + 15s read/write timeouts, 10s shutdown drain — keep handlers under that.
