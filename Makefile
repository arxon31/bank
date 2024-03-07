include .env
export $(shell sed 's/=.*//' .env)

.PHONY: build
build:
	go build -C cmd/app -o app

.PHONY: run
run: build
	./cmd/app/app

.PHONY: migrate-up
migrate-up:
	goose -dir migrations postgres $(PG_URL) up

.PHONY: migrate-down
migrate-down:
	goose -dir migrations postgres $(PG_URL) down
