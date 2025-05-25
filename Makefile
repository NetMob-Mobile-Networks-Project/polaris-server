.PHONY: sqlc db-up db-down polaris-server

sqlc:
	cd internal/storage && sqlc generate

db-up:
	./scripts/setup_db.sh

db-down:
	./scripts/setup_db.sh down

polaris-server:
	go build -o server ./cmd/polaris-server 