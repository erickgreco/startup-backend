PROTO_DIR := proto
PROTO_SRC := $(wildcard $(PROTO_DIR)/*.proto)
GO_OUT := .
COMPOSE_FILE := infra/development/docker/docker-compose.yml

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