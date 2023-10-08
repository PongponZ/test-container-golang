install:
	@go get ./...
	@go install github.com/cosmtrek/air@latest
	@go install github.com/vektra/mockery/v2@v2.20.0

dev:
	@air

test:
	@go test -v ./...

test-integration:
	@go test -v -race ./... -tags=integration