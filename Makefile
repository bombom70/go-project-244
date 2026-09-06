start:
	go run  cmd/gendiff/main.go

build:
	go build -o bin/gendiff ./cmd/gendiff

lint:
	golangci-lint run

fix:
	golangci-lint run --fix
