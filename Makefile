.PHONY: docker-compose-up
docker-compose-up:
	docker-compose up

.PHONY: docker-compose-down
docker-compose-down:
	docker-compose down

.PHONY: migration-up
migration-up:
	goose postgres "host=localhost user=postgres database=news_ai_feed password=postgres sslmode=disable" up

.PHONY: migration-down
migration-down:
	goose postgres "host=localhost user=postgres database=news_ai_feed password=postgres sslmode=disable" down

.PHONY: migration-status
migration-status:
	goose postgres "host=localhost user=postgres database=news_ai_feed password=postgres sslmode=disable" status