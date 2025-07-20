APP_ENV ?= development

MIGRATE_PATH := ./migrations

run:
	APP_ENV=$(APP_ENV) \
		go run cmd/main.go

migrate-up:
	@export $$(xargs < .env) && migrate \
		-path $(MIGRATE_PATH) \
		-database $$MIGRATE_URL up

migrate-down:
	@export $$(xargs < .env) && migrate \
		-path $(MIGRATE_PATH) \
		-database $$MIGRATE_URL down
