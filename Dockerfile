# Start from the official Golang image
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./cmd/polaris-server


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server .

COPY scripts/ ./scripts/
COPY internal/storage/schema/ ./internal/storage/schema/

ENV DB_USER=root
ENV DB_PASS=secret
ENV DB_HOST=localhost
ENV DB_PORT=3306
ENV DB_NAME=polaris

EXPOSE 8080

CMD ["./server"]