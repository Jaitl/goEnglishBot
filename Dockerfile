FROM golang:1.16 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o main ./app

FROM scratch

COPY --from=builder /app/main /app/main

ENTRYPOINT ["/app/main"]
