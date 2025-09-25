# --- Load .env kalau ada, dan export ke shell child ---
ifneq (,$(wildcard .env))
include .env
export
endif

# --- Konfigurasi umum ---
COMPOSE          ?= docker compose -f docker-compose.dev.yml
NETWORK          ?= boilerplate-template_go-network
MIGRATE_IMAGE    ?= migrate/migrate
MIGRATIONS_DIR   := $(PWD)/internal/database/migrations

# Fallback agar tetap jalan meski .env kosong
DB_HOST          ?= postgresdb
DB_PORT          ?= 5432
DB_USER          ?= postgres
DB_PASSWORD      ?= postgres
DB_NAME          ?= db_boilerplate

DB_URL           := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# Tunggu DB ready memakai pg_isready dari image postgres
WAIT_DB          := docker run --rm --network $(NETWORK) postgres:alpine \
	sh -c 'until pg_isready -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME); do echo "waiting for postgres..."; sleep 1; done'

# Default target
.DEFAULT_GOAL := start

# --- Daftar phony targets ---
.PHONY: start build lint gen \
        db-up wait-db \
        migration-% migrate-up migrate-down migrate-fresh \
        seed \
        docker-dev docker-prod docker-down docker-nuke docker-cache psql

# --- Go workflow ---
start:
	@go run cmd/api/main.go

build:
	@go build -o tmp/app ./cmd/api

lint:
	@golangci-lint run

# --- Compose / DB helpers ---
db-up:
	@$(COMPOSE) up -d postgresdb

wait-db:
	@$(WAIT_DB)

# --- Migration (pembuatan file) ---
# Contoh: make migration-create_users_table
# ":" akan diubah ke "_" (biar aman untuk nama file)
migration-%:
	@migrate create -ext sql -dir internal/database/migrations $(subst :,_,$*)

# --- Migration (apply via docker image 'migrate') ---
migrate-up: db-up wait-db
	@docker run --rm -v $(MIGRATIONS_DIR):/migrations --network $(NETWORK) \
		$(MIGRATE_IMAGE) -path=/migrations/ -database "$(DB_URL)" up

# Contoh:
#   make migrate-down step=2   → rollback 2 step
#   make migrate-down          → rollback semua

migrate-down: db-up wait-db
	@if [ -n "$(step)" ]; then \
		echo "⬇️  Migrating down $(step) step(s)..."; \
		docker run --rm -v $(MIGRATIONS_DIR):/migrations --network $(NETWORK) \
			$(MIGRATE_IMAGE) -path=/migrations/ -database "$(DB_URL)" down $(step); \
	else \
		echo "⬇️  Migrating down ALL steps..."; \
		docker run --rm -v $(MIGRATIONS_DIR):/migrations --network $(NETWORK) \
			$(MIGRATE_IMAGE) -path=/migrations/ -database "$(DB_URL)" down -all; \
	fi

migrate-fresh: migrate-down migrate-up
	@true

# Pakai: make migrate-force v=20250917120000
migrate-force:
	@docker run --rm -v $(MIGRATIONS_DIR):/migrations --network $(NETWORK) \
		$(MIGRATE_IMAGE) -path=/migrations/ -database "$(DB_URL)" force $(v)


# --- Seeder ---
seed: db-up wait-db
	@$(COMPOSE) run --rm app go run cmd/seed/main.go

# --- Docker orchestration convenience ---
docker-dev:
	@$(COMPOSE) up --build -d

docker-prod:
	@docker compose -f docker-compose.prod.yml up --build -d

docker-down:
	@$(COMPOSE) down --remove-orphans

# ⚠️ Akan menghapus container, images dan volumes.
docker-nuke:
	@$(COMPOSE) down --rmi all --volumes --remove-orphans

docker-cache:
	@docker builder prune -f

# --- PSQL shell ke DB di container ---
psql: db-up
	@$(COMPOSE) exec -it postgresdb psql -U $(DB_USER) -d $(DB_NAME)

# Single feature
# example: make gen feat=product-category

# Sub feature
# make gen feat=master/area
gen:
	@go run tools/gen.go $(feat)
# 	@goimports -w internal
