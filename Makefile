ifeq ($(POSTGRES_SETUP),)
	POSTGRES_SETUP_TEST := host=localhost user=postgres database=news_ai_feed password=postgres port=5432 sslmode=disable
endif

MIGRATION_FOLDER=$(CURDIR)/internal/storage/migrations

.PHONY: docker-compose-up docker-compose-down migration-up migration-down migration-status migration-create
docker-compose-up:
	docker-compose up

docker-compose-down:
	docker-compose down

migration-create:
	goose -dir "$(MIGRATION_FOLDER)" create "$(name)" sql

migration-up:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" up

migration-down:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" down

migration-status:
	goose -dir "$(MIGRATION_FOLDER)" postgres "$(POSTGRES_SETUP_TEST)" status