.PHONY: build
build:
	@go build -o ./dist/ ./cmd/library/main.go

.PHONY: run
run: build
	@./dist/main

.PHONY: test
test:
	@go test -v ./...
