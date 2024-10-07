# Name app
APP_NAME = server

dev:
	go run ./cmd/$(APP_NAME)

.PHONY: dev
