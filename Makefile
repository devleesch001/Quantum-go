APP_NAME := quantum-go

.PHONY: build run clean fmt tidy

build:
	@mkdir -p bin
	go build -v -o bin/$(APP_NAME) ./...

run:
	go run ./cmd/main.go

clean:
	rm -rf bin/

fmt:
	go fmt ./...

tidy:
	go mod tidy
