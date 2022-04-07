-include .env

IMAGE_NAME=line-messaging-app

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
	docker build --platform amd64 -f environment/server/Dockerfile -t ${IMAGE_NAME} .

.PHONY: run-prod
run-prod:
	docker run --rm -p 9000:80 ${IMAGE_NAME}

IMAGE_VERSION=latest

# cf. https://lightsail.aws.amazon.com/ls/docs/ja_jp/articles/amazon-lightsail-pushing-container-images
.PHONY: aws-push-image
aws-push-image: build-prod
	aws lightsail push-container-image --service-name line-messaging-app --label webhook --image ${IMAGE_NAME}:${IMAGE_VERSION}
