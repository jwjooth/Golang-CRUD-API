# Golang RESTful API

Clean-architecture RESTful API built with Go, Chi router, and GORM (MySQL). Implements CRUD operations for Products, Books, and Categories with strict layer boundaries, request validation, consistent JSON responses, pagination, Swagger documentation, and automated CI/CD.

---

## Table of Contents

- [Architecture](#architecture)
  - [Overview](#overview)
  - [Architecture Flow Diagram](#architecture-flow-diagram)
  - [Layer Boundaries and Rules](#layer-boundaries-and-rules)
  - [Project Structure](#project-structure)
- [Requirements](#requirements)
- [Installation](#installation)
- [Environment Variables](#environment-variables)
- [Database Setup](#database-setup)
- [Running Locally](#running-locally)
- [Running Tests](#running-tests)
- [Swagger Documentation](#swagger-documentation)
- [Docker](#docker)
- [Deployment](#deployment)
- [API Reference & Examples](#api-reference--examples)
  - [Response Envelope](#response-envelope)
  - [Health Check](#health-check)
  - [Products Endpoints](#products-endpoints)
  - [Books Endpoints](#books-endpoints)
  - [Categories Endpoints](#categories-endpoints)
- [License](#license)

---

## Architecture

### Overview

The application follows the **Clean Architecture** pattern to achieve separation of concerns, high maintainability, and testability. Incoming HTTP requests pass strictly through layers from routing down to persistence:

### Architecture Flow Diagram

```
Client
  ↓
Router (Chi Router + Middleware)
  ↓
Controller (Decode parameters/body, invoke service, encode JSON)
  ↓
Service (Validation, business rules, entity ↔ DTO mapping)
  ↓
Repository (Data access with GORM, returns entities and error sentinels)
  ↓
MySQL Database
```

### Layer Boundaries and Rules

| Layer | Package | Depends on | Must NOT contain |
|---|---|---|---|
| **HTTP** | `controller/` | `service`, `payload`, `helper` | SQL queries, business logic, direct GORM access |
| **Business** | `service/` | `repository`, `payload`, `entities`, `helper` | HTTP-specific code (`net/http`), direct SQL queries |
| **Persistence** | `repository/` | `entities` | DTOs, HTTP status codes, validation logic |
| **DTOs** | `payload/` | `entities` (mappers only) | Database logic, HTTP handling |
| **Shared** | `helper/`, `config/` | Standard library, drivers | Core business/domain logic |

### Project Structure

```
.
├── main.go                     # Application entry point, dependency wiring, routes, graceful shutdown
├── config/
│   └── config.go               # Configuration loader, DB connection pool, auto-migration
├── controller/                 # HTTP handlers (parse input, invoke service, write envelope)
│   ├── BookController.go
│   ├── CategoryController.go
│   └── ProductController.go
├── service/                    # Business logic, payload validation, DTO mapping
│   ├── BookService.go
│   ├── CategoryService.go
│   └── ProductService.go
├── repository/                 # Data persistence layer using GORM
│   ├── BookRepository.go
│   ├── CategoryRepository.go
│   └── ProductRepository.go
├── entities/                   # GORM database models
│   ├── BookEntity.go
│   ├── CategoryEntity.go
│   └── ProductEntity.go
├── payload/                    # Request/response DTOs and mappers
│   ├── BookPayload.go
│   ├── CategoryPayload.go
│   └── ProductPayload.go
├── helper/                     # Shared helpers: response envelopes, error structures, query parsers
│   ├── error.go
│   └── response.go
├── docs/                       # Swagger 2.0 generated documentation (swag init)
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── test/                       # Unit and integration test suites
│   ├── controller/
│   ├── repository/
│   └── service/
├── .github/workflows/          # GitHub Actions CI/CD workflows
│   └── go-ci.yml
├── http-client.http            # HTTP requests for testing with JetBrains / VS Code REST clients
├── .env.example                # Sample environment variables
├── go.mod / go.sum             # Go module dependencies
└── README.md                   # Project documentation
```

---

## Requirements

- **Go**: 1.25 or higher (tested with Go 1.26+)
- **MySQL**: 8.0 or higher
- **Git**: For cloning the repository
- **curl** or **REST Client** (Postman, Thunder Client, etc.) for testing endpoints
- *(Optional)* **Docker & Docker Compose**: For containerized database or app setup
- *(Optional)* **swag CLI**: `github.com/swaggo/swag/cmd/swag` for generating Swagger specifications

---

## Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/jwjooth/Product-CRUD-API.git
   cd golang-restful-api
   ```

2. **Download dependencies:**
   ```bash
   go mod download
   # or ensure module dependencies are tidy:
   go mod tidy
   ```

---

## Environment Variables

The application can read environment variables from a `.env` file in the project root or from system environment variables. The `.env` file is optional; defaults are used if values are not specified.

Copy `.env.example` to `.env`:
```bash
cp .env.example .env
```

### Available Variables

| Variable | Type | Default | Description |
|---|---|---|---|
| `APP_PORT` | `string` | `6767` | Port on which the HTTP server will listen. |
| `APP_URL` | `string` | *(empty)* | Optional base URL for Swagger documentation host & scheme (e.g. `https://api.example.com`). |
| `DB_HOST` | `string` | `localhost` | MySQL host address. |
| `DB_PORT` | `string` | `3306` | MySQL port. |
| `DB_USER` | `string` | `root` | MySQL user. |
| `DB_PASSWORD` | `string` | *(empty)* | MySQL user password. |
| `DB_NAME` | `string` | `products` | MySQL database name. |
| `DB_MAX_IDLE_CONNS` | `int` | `10` | Maximum number of idle connections in pool. |
| `DB_MAX_OPEN_CONNS` | `int` | `100` | Maximum number of open connections in pool. |

---

## Database Setup

1. **Create the MySQL Database:**
   Connect to MySQL and create the database if it doesn't already exist:
   ```sql
   CREATE DATABASE IF NOT EXISTS products CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   ```
   Or via command-line:
   ```bash
   mysql -h localhost -P 3306 -u root -p -e "CREATE DATABASE IF NOT EXISTS products CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
   ```

2. **Auto Migration:**
   When the server boots up, `config.Migrate` executes GORM's `AutoMigrate` to verify and create the necessary tables.
   > [!NOTE]
   > For production deployments, consider utilizing versioned migration tools such as `golang-migrate`.

---

## Running Locally

1. **Ensure MySQL is running** and credentials in `.env` match your database.
2. **Start the application:**
   ```bash
   go run .
   ```
3. **Verify the server is running:**
   ```bash
   curl -i http://localhost:6767/healthz
   ```
   Expected response:
   ```json
   {
     "message": "ok",
     "data": {
       "status": "up"
     }
   }
   ```

4. **Build a binary locally:**
   ```bash
   # Verify code and build binary
   go vet ./...
   go build -o api .
   
   # Run binary
   ./api
   ```

---

## Running Tests

Run the test suite across all packages:
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage summary
go test -cover ./...

# Run specific package tests
go test ./test/service/... -v
```

CI workflows also run automatically on GitHub Actions for pull requests and pushes to `main` (see `.github/workflows/go-ci.yml`).

---

## Swagger Documentation

Interactive Swagger API documentation is pre-generated and accessible via browser.

- **Swagger UI URL:** [http://localhost:6767/swagger/index.html](http://localhost:6767/swagger/index.html)

### Regenerating Swagger Documentation

If you update controller Swagger annotations or payload structures:
1. Install `swag` CLI if not already installed:
   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   ```
2. Re-generate Swagger files:
   ```bash
   swag init -g main.go -o docs
   ```

---

## Docker

### Running MySQL via Docker

If you do not have MySQL installed locally, run a MySQL 8 container:
```bash
docker run --name mysql-restful-api \
  -e MYSQL_ROOT_PASSWORD=secret \
  -e MYSQL_DATABASE=products \
  -p 3306:3306 \
  -d mysql:8.0
```

### Dockerfile for API Application

Create a `Dockerfile` for containerizing the Go application:

```dockerfile
# Multi-stage build
FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o api .

# Final minimal image
FROM alpine:3.20

WORKDIR /app
RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/api /app/api

EXPOSE 6767

ENTRYPOINT ["/app/api"]
```

Build and run the Docker image:
```bash
# Build Docker image
docker build -t golang-restful-api:latest .

# Run Docker container linking to host database
docker run -p 6767:6767 \
  -e DB_HOST=host.docker.internal \
  -e DB_USER=root \
  -e DB_PASSWORD=secret \
  -e DB_NAME=products \
  golang-restful-api:latest
```

---

## Deployment

The application compiles into a single, self-contained binary that can be deployed on any server, VM, or PaaS platform (e.g. Railway, Render, Fly.io, AWS EC2 / ECS, Kubernetes).

### Production Considerations

1. **Environment Variables**: Never commit `.env` to version control. Set environment variables directly in your hosting platform dashboard.
2. **Reverse Proxy & TLS**: Place the API behind Nginx, Caddy, Cloudflare, or an API gateway for HTTPS termination and rate limiting.
3. **Database Pool Sizing**: Adjust `DB_MAX_OPEN_CONNS` and `DB_MAX_IDLE_CONNS` according to your database tier `max_connections` and instance scaling count.
4. **App URL for Swagger**: Provide `APP_URL` (e.g., `APP_URL=https://your-domain.com`) so the Swagger UI can resolve the correct host and scheme for "Try it out" calls.
5. **Graceful Shutdown**: The app handles `SIGINT` and `SIGTERM` signals with a 10-second drain window before terminating active connections.

---

## API Reference & Examples

**Base URL**: `http://localhost:6767/api/v1`

### Response Envelope

All API endpoints return consistent JSON envelopes:

#### Success Response
```json
{
  "message": "success",
  "data": { ... }
}
```

#### Paginated Success Response
```json
{
  "message": "success",
  "data": [ ... ],
  "meta": {
    "page": 1,
    "per_page": 10,
    "total": 45,
    "total_pages": 5
  }
}
```

#### Error Response
```json
{
  "error": "description of the error"
}
```

---

### Health Check

- **Endpoint**: `GET /healthz`
- **Description**: Verifies if the service is running.

```bash
curl -s http://localhost:6767/healthz
```

---

### Products Endpoints

#### Summary Table
| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/v1/products?page=1&per_page=10` | List products with pagination |
| `POST` | `/api/v1/products` | Create a new product |
| `GET` | `/api/v1/products/{id}` | Get product by ID |
| `PUT` | `/api/v1/products/{id}` | Update product by ID |
| `DELETE` | `/api/v1/products/{id}` | Delete product by ID |

#### Validation Rules (Product)
- `name`: Required, 1–255 characters.
- `description`: Optional, max 5000 characters.
- `price`: Required, greater than `0`.
- `stock`: Required / default `0`, greater than or equal to `0`.

#### Examples

**1. List Products (Paginated)**
```bash
curl -s "http://localhost:6767/api/v1/products?page=1&per_page=5"
```

**2. Create Product**
```bash
curl -s -X POST "http://localhost:6767/api/v1/products" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Mechanical Keyboard RGB",
    "description": "Tenkeyless mechanical gaming keyboard with red switches",
    "price": 89.99,
    "stock": 35
  }'
```

**3. Get Product by ID**
```bash
curl -s "http://localhost:6767/api/v1/products/1"
```

**4. Update Product**
```bash
curl -s -X PUT "http://localhost:6767/api/v1/products/1" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Mechanical Keyboard RGB v2",
    "description": "Updated switches and wireless connectivity",
    "price": 99.99,
    "stock": 20
  }'
```

**5. Delete Product**
```bash
curl -s -X DELETE "http://localhost:6767/api/v1/products/1"
```

---

### Books Endpoints

#### Summary Table
| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/v1/books?page=1&per_page=10` | List books with pagination |
| `POST` | `/api/v1/books` | Create a new book |
| `GET` | `/api/v1/books/{id}` | Get book by ID |
| `PUT` | `/api/v1/books/{id}` | Update book by ID |
| `DELETE` | `/api/v1/books/{id}` | Delete book by ID |

#### Validation Rules (Book)
- `title`: Required.
- `category_id`: Required.
- `author`: Required.
- `stock`: Required.

#### Examples

**1. List Books**
```bash
curl -s "http://localhost:6767/api/v1/books?page=1&per_page=10"
```

**2. Create Book**
```bash
curl -s -X POST "http://localhost:6767/api/v1/books" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Clean Architecture: A Craftsman Guide to Software Structure and Design",
    "category_id": 1,
    "author": "Robert C. Martin",
    "stock": 15
  }'
```

**3. Get Book by ID**
```bash
curl -s "http://localhost:6767/api/v1/books/1"
```

**4. Update Book**
```bash
curl -s -X PUT "http://localhost:6767/api/v1/books/1" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Clean Architecture (2nd Edition)",
    "category_id": 1,
    "author": "Robert C. Martin",
    "stock": 25
  }'
```

**5. Delete Book**
```bash
curl -s -X DELETE "http://localhost:6767/api/v1/books/1"
```

---

### Categories Endpoints

#### Summary Table
| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/v1/categories` | List all categories (unpaginated) |
| `POST` | `/api/v1/categories` | Create a new category |
| `GET` | `/api/v1/categories/{id}` | Get category by ID |
| `PUT` | `/api/v1/categories/{id}` | Update category by ID |
| `DELETE` | `/api/v1/categories/{id}` | Delete category by ID |

#### Examples

**1. List All Categories**
```bash
curl -s "http://localhost:6767/api/v1/categories"
```

**2. Create Category**
```bash
curl -s -X POST "http://localhost:6767/api/v1/categories" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Technology & Programming"
  }'
```

**3. Get Category by ID**
```bash
curl -s "http://localhost:6767/api/v1/categories/1"
```

**4. Update Category**
```bash
curl -s -X PUT "http://localhost:6767/api/v1/categories/1" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Computer Science & Software Engineering"
  }'
```

**5. Delete Category**
```bash
curl -s -X DELETE "http://localhost:6767/api/v1/categories/1"
```

---

## License

This project is licensed under the [MIT License](LICENSE) (or refer to repository settings for distribution).
