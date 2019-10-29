FROM golang:1.13 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o main ./app

FROM ubuntu:18.04

RUN apt-get update && apt-get --no-install-recommends install -y ca-certificates ffmpeg

COPY --from=builder /app/main /app/main

RUN useradd appuser
RUN mkdir /tmp_data
RUN chown appuser: /tmp_data

USER appuser

ENTRYPOINT ["/app/main"]
