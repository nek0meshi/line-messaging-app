.PHONY: up
up:
	docker compose up

.PHONY: down
down:
	docker compose down

.PHONY: sh
sh:
	docker compose run --rm server bash
