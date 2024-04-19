run: build
	@./bin/wapp

build:
	@go build -o bin/wapp cmd/main.go