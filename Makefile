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

migrate-up:
	migrate -source file://internal/user/infras/mysql/migrations -database "$$(grep DATABASE_DSN cmd/user/.env | cut -d '=' -f 2-)" up

migrate-down:
	migrate -source file://internal/user/infras/mysql/migrations -database "$$(grep DATABASE_DSN cmd/user/.env | cut -d '=' -f 2-)" down

docker-build:
	docker build -t user:latest -f docker/Dockerfile-user .

docker-run:
	docker run -d -p 8080:8080 --env-file cmd/user/.env user:latest

clean:
	rm -rf bin/*