-include infra/development/docker/.env
export

PROTO_DIR := proto
PROTO_SRC := $(wildcard $(PROTO_DIR)/*.proto)
GO_OUT := .
COMPOSE_FILE := infra/development/docker/docker-compose.yml

MIGRATIONS_DIR := services/crm-pm-service/internal/infrastructure/db/migrations
DB_DRIVER := postgres
DB_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

.PHONY: generate-proto
generate-proto:
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(GO_OUT) \
		--go-grpc_out=$(GO_OUT) \
		$(PROTO_SRC)

.PHONY: db-up db-down db-logs

db-up: 
	docker compose -f $(COMPOSE_FILE) up -d

db-down: 
	docker compose -f $(COMPOSE_FILE) down

db-logs: 
	docker compose -f $(COMPOSE_FILE) logs -f

.PHONY: goose-up goose-down goose-status goose-create
goose-up:
	goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(DB_URL)" up

goose-down:
	goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(DB_URL)" down

goose-status:
	goose -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(DB_URL)" status

goose-create:
	@read -p "Migration name: " name; \
	goose -dir $(MIGRATIONS_DIR) create $$name sql