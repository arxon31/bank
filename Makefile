include .env

.PHONY: build
build:
	go build -C cmd/app -o app

.PHONY: run
run: build
	 ./cmd/app/app

.PHONY: migrate-up
migrate-up:
	migrate -path migrations -database '$(PG_URL)?sslmode=disable' up
