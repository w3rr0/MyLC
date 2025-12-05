FROM golang:1.25 AS builder

LABEL authors="w3rr0"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o server ./cmd/server/main.go

FROM ubuntu:24.10
WORKDIR /app
COPY --from=builder /app/server .
ENV DATABASE_URL="postgres://postgres:password@db:5432/mylc?sslmode=disable"
CMD ["./server"]