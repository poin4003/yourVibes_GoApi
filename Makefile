# Name app
APP_NAME = server

dev:
	go run ./cmd/$(APP_NAME)

swag:
	swag init -g ./cmd/server/main.go -o ./cmd/swag/docs

.PHONY: dev
