# Golang Boilerplate (golang-boilerplate)

Boilerplate RESTful API, built with **Go, Fiber, GORM** and **PostgreSQL**.

---

## 📦 Tech Stack

- **Go + Fiber** — Web framework
- **GORM** — ORM for PostgreSQL
- **PostgreSQL** — Relational database
- **go-playground/validator** — Input validation
- **JWT** — Authentication
- **Logrus** — Logging
- **Fiber middleware** — Rate limiting, CORS, recovery, logger
- **Air** — Hot reload for development
- **Docker + Docker Compose** — Containerization

---

## 🚀 Getting Started

### 1. Clone Project

```bash
git clone https://github.com/hafizhproject45/Golang-Boilerplate.git
cd golang-boilerplate
```

### 2. Install Dependencies

```bash
go mod tidy
```

### 3. Configure Environment

Copy .env.example to .env and adjust the variables (e.g. DATABASE_URL, JWT secrets, etc).

```bash
cp .env.example .env
```

### 5. Setup Docker

Run initial docker.

```bash
make docker-dev
```

### 4. Migrate Database

Run initial migrations and generate views.

```bash
make migrate-up
```

### 5. Run App

Run project via Docker

### 6. Create New Module

```bash
make gen feat=user
```

output:

```bash
cmd/
├── api/
│   └── main.go            # Application entrypoint (initialize Fiber, load config, connect DB, register route)

internal/
├── config/                # App config (env loader, logger, app settings)
│
├── database/              # Database connection + migration setup
│
├── middleware/            # Global Fiber middleware (auth, logger, recovery, rate limiting)
│
├── modules/               # Feature modules (users, products, suppliers, etc.)
│   ├── <module>/
│   │   ├── controllers/   # HTTP handler layer (receive request, call service, return response)
│   │   ├── dto/           # Data Transfer Objects (request & response payloads, separate from models)
│   │   ├── models/        # GORM models (represent database tables/entities)
│   │   ├── repositories/  # Data access layer (queries to DB, CRUD abstraction)
│   │   ├── services/      # Business logic layer (process rules, orchestrate repository calls)
│   │   ├── validation/    # Request validation (custom rules per module)
│   │   ├── module.go      # Module bootstrapper (wire controller, service, repository together)
│   │   └── route.go       # Module route (register module routes into Fiber app)
│
├── repository/            # Shared repositories (reusable DB access layer across multiple modules)
│
├── response/              # Standardized API responses (success, error, pagination)
│
├── utils/                 # Helper functions (JWT, hashing, constants, enums, etc.)
│
├── validation/            # Shared request validation structs & rules
│
└── route/                 # Central route aggregator (load all module routes into main app)
```

## ✨ Author

IT Development PT Mitra Berlian Unggas

## 📃 License

This project is private. All rights reserved.
