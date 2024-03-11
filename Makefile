include .env
export $(shell sed 's/=.*//' .env)

.PHONY: build
build:
	go build -C cmd/app -o app

.PHONY: run-publisher
run-publisher:
	go run cmd/publisher/main.go

.PHONY: run
run: build migrate-up
	./cmd/app/app

.PHONY: migrate-up
migrate-up:
	goose -dir migrations postgres $(PG_URL) up

.PHONY: migrate-down
migrate-down:
	goose -dir migrations postgres $(PG_URL) down

.PHONY: docker-build
docker-build:
	docker build . -t $(APP_NAME):latest
	docker build . -t producer:latest -f Dockerfile-producer

.PHONY: compose-up
compose-up:
	docker-compose up -d

.PHONY: test
test:
	go test -v -cover  ./...

.PHONY: lint
lint:
	golangci-lint run --fast --timeout 1m --config=./.golangci.yml

.PHONY: compose-app
compose-app: docker-build compose-up