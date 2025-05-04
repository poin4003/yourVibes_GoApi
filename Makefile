# Name app
APP_NAME = server

# Config by OS
ifeq ($(OS),Windows_NT)
	SET_ENV = set
else
	SET_ENV = export
endif

# Config command
dev:
	@echo "Running in development mode"
	@$(SET_ENV) YOURVIBES_SERVER_CONFIG_FILE=dev&&go run ./cmd/$(APP_NAME)

prod:
	@echo "Running in production mode"
	@$(SET_ENV) YOURVIBES_SERVER_CONFIG_FILE=prod&&go run ./cmd/$(APP_NAME)

cloud:
	@echo "Running in cloud mode"
	@$(SET_ENV) YOURVIBES_SERVER_CONFIG_FILE=cloud&&go run ./cmd/$(APP_NAME)

migrate:
	@echo "Running migrations with config: $(CONFIG_FILE)"
	@$(SET_ENV) YOURVIBES_SERVER_CONFIG_FILE=$(CONFIG_FILE)&&go run ./cmd/cli/postgresql/migrate.go

swag:
	swag init -g ./cmd/server/main.go -o ./cmd/swag/docs

# gRPC code generation
PROTO_DIR := ./cmd/proto
OUT_DIR := ./internal/infrastructure/pkg/grpc

gen-grpc:
	@echo Generating gRPC code for $(FILE).proto...
	protoc -I=$(PROTO_DIR) \
		--go_out=$(OUT_DIR) \
		--go-grpc_out=$(OUT_DIR) \
		$(PROTO_DIR)/$(FILE).proto


docker_build:
	docker-compose up -d --build

docker_stop:
	docker-compose down

docker_up:
	docker-compose up -d

.PHONY: dev prod cloud migrate swag gen-grpc docker_build docker_stop docker_up
