docker:
	docker build -t go-english-bot .

lint:
	golangci-lint run ./...

test:
	go test ./...
