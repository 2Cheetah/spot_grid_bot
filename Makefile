.DEFAULT_GOAL := build

.PHONY:fmt vet build run
fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	go build ./cmd/echo_path/...

clean:
	go clean

run:
	set -a; source .env; go run ./cmd/client/main.go; set +a
