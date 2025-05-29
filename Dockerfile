FROM golang:1.24-alpine

WORKDIR /app

RUN apk add --no-cache mysql-client

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./cmd/polaris-server

RUN mkdir -p /app/scripts
COPY scripts/setup_db.sh /app/scripts/
COPY scripts/down_db.sh /app/scripts/
RUN chmod +x /app/scripts/*.sh

EXPOSE 8080

RUN echo '#!/bin/sh' > /app/entrypoint.sh && \
    echo 'cd /app' >> /app/entrypoint.sh && \
    echo 'sh /app/scripts/setup_db.sh' >> /app/entrypoint.sh && \
    echo 'exec ./server' >> /app/entrypoint.sh

RUN chmod +x /app/entrypoint.sh

ENTRYPOINT ["sh", "/app/entrypoint.sh"]
