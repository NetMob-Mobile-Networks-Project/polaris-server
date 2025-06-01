.PHONY: sqlc db-up db-down polaris-server

docker-compose-up:
	cp .env.example .env
	docker compose up --build --force-recreate -d

docker-compose-down:
	docker compose down	-v

sqlc:
	cd internal/storage && sqlc generate

db-up:
	./scripts/setup_db.sh

db-down:
	./scripts/down_db.sh

polaris-server:
	go build -o server ./cmd/polaris-server 