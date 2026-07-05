.PHONY: build run test clean

build:
	go build -o tmp/main.exe ./cmd/api

run:
	go run ./cmd/api/main.go

test:
	go test ./...

clean:
	go clean
