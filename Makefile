.PHONY: build
build:
	@go build -o ./dist/ ./cmd/main.go

.PHONY: build
run: build
	@./dist/main

.PHONY: test
test:
	@go test -v ./...