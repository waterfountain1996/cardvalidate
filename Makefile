format:
	@gofmt -w -l .

test:
	@go test ./...
