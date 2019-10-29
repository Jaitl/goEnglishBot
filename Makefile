.PHONY: test bin

build:
	go mod download
	go build -o main ./app

docker:
	docker build -t go-english-bot .

lint:
	golangci-lint run ./...

test:
	go test ./...
