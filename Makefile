BINARY_NAME=main.out

start: clean dev
start_test: clean test

dev:
	go run cmd/main.go

build:
	go build cmd/main.go

clean:
	go clean

test:
	go test -v ./...