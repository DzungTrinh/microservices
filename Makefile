.PHONY: all build run test sqlc migrate-up migrate-down docker-build docker-run clean

all: build

build:
	go build -o bin/user ./cmd/user

run:
	go run ./cmd/user

test:
	go test ./internal/user/...

sqlc:
	sqlc generate

migrate-create:
	migrate create -ext sql -dir db/migrations/user -seq create_user_service

docker-build:
	docker build -t user:latest -f docker/Dockerfile-user .

docker-run:
	docker run -d -p 8080:8080 --env-file cmd/user/.env user:latest

docker-compose-up:
	docker-compose up

docker-compose-up-build:
	docker-compose up -d --build

docker-compose-down:
	docker-compose down

clean:
	rm -rf bin/*
