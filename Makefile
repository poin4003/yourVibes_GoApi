# Name app
APP_NAME = server

dev:
	go run ./cmd/$(APP_NAME)

migrate:
	go run ./cmd/cli/postgresql/migrate.go

swag:
	swag init -g ./cmd/server/main.go -o ./cmd/swag/docs

docker_build:
	docker-compose up -d --build

docker_stop:
	docker-compose down

docker_up:
	docker-compose up -d

.PHONY: dev
