run:
	go run cmd/main.go

include .env
export

MIGRATE_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)
MIGRATE_PATH=./migrations

migrate-up:
	migrate -path $(MIGRATE_PATH) -database "$(MIGRATE_URL)" up

migrate-down:
	migrate -path $(MIGRATE_PATH) -database "$(MIGRATE_URL)" down
