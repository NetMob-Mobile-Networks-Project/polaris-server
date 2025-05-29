FROM golang:1.24-alpine AS builder

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

EXPOSE 8080


RUN echo '#!/bin/sh' > /app/entrypoint.sh && \
    echo '/app/scripts/setup_db.sh' >> /app/entrypoint.sh && \
    echo 'exec ./server' >> /app/entrypoint.sh

RUN chmod +x /app/entrypoint.sh

ENTRYPOINT ["sh", "/home/app/entrypoint.sh"]
