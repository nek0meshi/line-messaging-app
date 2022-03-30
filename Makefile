.PHONY: up
up:
	docker compose up

.PHONY: down
down:
	docker compose down

.PHONY: sh
sh:
	docker compose run --rm server bash

.PHONY: build-prod
build-prod:
	docker build -f environment/server/Dockerfile -t nek0meshi/line-messaging-app .

.PHONY: run-prod
run-prod:
	docker run --rm nek0meshi/line-messaging-app
