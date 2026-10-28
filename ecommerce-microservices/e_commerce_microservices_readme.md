# Resilient E-Commerce Microservices Platform

A production-ready, distributed e-commerce backend built with **Go**, **Gin**, **GORM**, **PostgreSQL**, **GraphQL**, and **HashiCorp Consul**. The system implements clean architecture across separate microservices with an API Gateway acting as a unified entry point, featuring stateless JWT/RBAC security, circuit breaker resilience, exponential backoff, structured JSON logging (`slog`), and automated database migrations.

## 🛠️ Tools, Frameworks & Technologies Used

### Core & Frameworks
* **Go (Golang 1.25+)**: Core programming language.
* **Gin Web Framework (`github.com/gin-gonic/gin`)**: High-performance HTTP web framework routing the gateway and microservices.
* **GORM (`gorm.io/gorm`)**: Object Relational Mapping library for PostgreSQL data access and automated schema migrations (`AutoMigrate`).
* **GORM PostgreSQL Driver (`gorm.io/driver/postgres`)**: PostgreSQL driver for GORM built on `pgx/v5`.
* **GraphQL Go (`github.com/graphql-go/graphql` & `github.com/graphql-go/handler`)**: Schema declaration, query execution, mutations, and embedded GraphQL Playground / GraphiQL interface.

### Infrastructure & Discovery
* **PostgreSQL (15+)**: Primary relational database storing product catalogs and order transactions.
* **HashiCorp Consul**: Centralized service registry, health-checking agent, and dynamic endpoint catalog.
* **PowerShell**: Automation script (`register-consul.ps1`) for direct HTTP REST registration with Consul.

### Security & Resilience
* **JWT (`github.com/golang-jwt/jwt/v5`)**: Stateless HMAC-SHA256 token generation, claims verification, and expiration enforcement.
* **RBAC (Role-Based Access Control)**: Granular permission enforcement (`ADMIN` vs. `USER`) at the gateway layer.
* **Custom Circuit Breaker State Machine**: `CLOSED`, `OPEN`, and `HALF_OPEN` state protection against cascading downstream outages.
* **Context-Aware Exponential Backoff**: Non-blocking retry mechanism adhering to `context.Context` cancellation deadlines.

### Observability & Tooling
* **Structured Logging (`log/slog`)**: JSON-formatted logging outputting request timings, status codes, and IP addresses.
* **Distributed Tracing Header**: `X-Request-ID` generation and context propagation across service boundaries.
* **Godotenv (`github.com/joho/godotenv`)**: Environment variable loader.
* **AWS Elastic Beanstalk & CodeBuild**: `Procfile` and `buildspec.yml` deployment definitions.
* **Testing Tools**: Standard Go testing suite (`go test`), SoapUI, and Postman for REST / GraphQL API validation.

## 📁 Project File Structure

```
ecommerce-microservices/
├── .env
├── .env.example
├── .gitignore
├── buildspec.yml
├── examples.graphql
├── go.mod
├── go.sum
├── Procfile
├── register-consul.ps1
├── schema.sql
│
├── database/
│   └── postgres.go
│
├── gateway/
│   ├── main.go
│   ├── client/
│   │   └── consul_client.go
│   ├── config/
│   │   └── config.go
│   └── graphql/
│       ├── handler.go
│       └── schema.go
│
├── pkg/
│   ├── auth/
│   │   └── jwt.go
│   ├── middleware/
│   │   ├── logging.go
│   │   ├── recovery.go
│   │   └── request_id.go
│   └── resilience/
│       └── circuit_breaker.go
│
├── services/
│   ├── order/
│   │   ├── main.go
│   │   ├── config/
│   │   │   └── config.go
│   │   ├── controller/
│   │   │   └── order_controller.go
│   │   ├── model/
│   │   │   └── order.go
│   │   ├── repository/
│   │   │   └── order_repository.go
│   │   └── service/
│   │       └── order_service.go
│   │
│   └── product/
│       ├── main.go
│       ├── config/
│       │   └── config.go
│       ├── controller/
│       │   └── product_controller.go
│       ├── model/
│       │   └── product.go
│       ├── repository/
│       │   └── product_repository.go
│       └── service/
│           └── product_service.go
│
└── tests/
    ├── jwt_unit_test.go
    ├── order_unit_test.go
    └── product_unit_test.go
```

## 🚀 Running the Services

Open 3 separate terminals from the project root:

### 1. Start Product Service
```bash
cd services/product
go run .
```

### 2. Start Order Service
```bash
cd services/order
go run .
```

### 3. Start API Gateway
```bash
cd gateway
go run .
```

*Access GraphQL Playground at `http://localhost:8080/playground`.*

## 🧪 Testing

Run all unit and concurrency tests from the repository root:
```bash
go test -v ./tests/...
```