.PHONY: test bin

build:
	go mod download
	go build -o main ./app

docker:
	docker build -t go-english-bot .

lint:
	golangci-lint run ./...

lint-ci:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint
	make lint

test:
	go test ./...
