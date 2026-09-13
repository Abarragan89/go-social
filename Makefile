include .env

MIGRATION_PATH = ./cmd/migrate/migrations

# EXAMPLE: make migration <migration-name>
.PHONY: migrate-create
migration:
	@migrate create -seq -ext sql -dir $(MIGRATION_PATH) $(filter-out $@,$(MAKECMDGOALS))

.PHONY: migrate-up
migrate-up:
	@migrate -path=$(MIGRATION_PATH) -database=$(DB_ADDR) up

.PHONY: migrate-down
migrate-down:
	@migrate -path=$(MIGRATION_PATH) -database=$(DB_ADDR) down $(filter-out $@,$(MAKECMDGOALS))

.PHONY:	migrate-drop:
migrate-drop:
	@migrate -path=./cmd/migrate/migrations -database=$(DB_ADDR) drop -f