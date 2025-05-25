.PHONY: sqlc db-up db-down

sqlc:
	cd internal/storage && sqlc generate

db-up:
	./scripts/setup_db.sh

db-down:
	./scripts/setup_db.sh down 