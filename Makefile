APP_NAME := quantum-go

.PHONY: build run clean fmt tidy

build:
	@mkdir -p bin
	go build -o bin/$(APP_NAME) ./cmd/main.go

run:
	go run ./cmd/$(APP_NAME)

clean:
	rm -rf bin/

fmt:
	go fmt ./...

tidy:
	go mod tidy
