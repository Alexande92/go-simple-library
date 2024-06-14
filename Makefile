.PHONY: build
build:
	@go build -o ./dist/ ./cmd/library/main.go

.PHONY: run
run: build
	@./dist/main

.PHONY: test
test:
	@go test -v ./internal/... -coverprofile=coverage.out -covermode atomic
	@go tool cover -html=coverage.out -o coverage.html
	@go tool cover -func=coverage.out
