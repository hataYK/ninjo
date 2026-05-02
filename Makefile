.PHONY: up down build logs generate-ent generate-swagger generate-api-client generate-all

# Docker
up:
	docker compose up -d

down:
	docker compose down

build:
	docker compose build

logs:
	docker compose logs -f

logs-be:
	docker compose logs -f backend

logs-fe:
	docker compose logs -f frontend

# Code generation
generate-ent:
	cd backend && go generate ./ent

generate-swagger:
	cd backend && swag init -g cmd/server/main.go -o docs/

generate-api-client:
	cd frontend && npx orval

generate-all: generate-ent generate-swagger generate-api-client

# DB
shell-db:
	docker compose exec db psql -U ninjo

shell-be:
	docker compose exec backend sh

shell-fe:
	docker compose exec frontend sh
