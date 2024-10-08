DOCKER_IMAGE_TAG ?= cardvalidate

all: build

run: build
	@./bin/server

build:
	@go build -o ./bin/server ./cmd/server

docker-build:
	@docker build -t $(DOCKER_IMAGE_TAG) .

format:
	@gofmt -w -l .

test:
	@go test ./...
