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
	docker build -f environment/server/Dockerfile -t ${IMAGE_NAME} .

.PHONY: run-prod
run-prod:
	docker run --rm ${IMAGE_NAME}

.PHONY: aws-login
aws-login:
	aws ecr get-login-password | \
		docker login --username AWS --password-stdin ${ECR_URL}

IMAGE_VERSION=latest

# cf. https://docs.aws.amazon.com/ja_jp/AmazonECR/latest/userguide/repository-policy-examples.html
.PHONY: aws-push-image
aws-push-image:
	docker tag ${IMAGE_NAME}:${IMAGE_VERSION} ${ECR_URL}/${IMAGE_NAME}:${IMAGE_VERSION}
	docker push ${ECR_URL}/${IMAGE_NAME}:${IMAGE_VERSION}
