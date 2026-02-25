include .env

migration:
	goose -dir internal/db/migrations create $(name) sql

migrate:
	goose -dir 'internal/db/migrations' -allow-missing  postgres ${DATABASE_URL} up

migrate-up-one:
	goose -dir 'internal/db/migrations' postgres ${DATABASE_URL} up-by-one

migratedbs:
	make migrate

rollback:
	goose -dir 'internal/db/migrations' postgres ${DATABASE_URL} down

resetdb:
	goose -dir 'internal/db/migrations' postgres ${DATABASE_URL} reset

resetmigrations:
	make resetdb

status:
	goose -dir 'internal/db/migrations' -allow-missing  postgres ${DATABASE_URL} status

api:
	go run cmd/api/*.go

# Clear Go build cache so the latest source is used (fixes stale/cached build errors)
clean-cache:
	go clean -cache
	go clean -testcache

# Rebuild from scratch after clearing caches
rebuild: clean-cache
	go build ./...

build:
	docker-compose up --build -d