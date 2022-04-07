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

.PHONY: fmt
fmt:
	docker compose run --rm server go fmt

IMAGE_VERSION=latest

LIGHTSAIL_SERVICE_NAME=line-messaging-app
LIGHTSAIL_CONTAINER_NAME=line-messaging-app
LIGHTSAIL_IMAGE_NAME_PREFIX=:line-messaging-app.webhook.
LIGHTSAIL_IMAGE_VERSION=

# cf. https://lightsail.aws.amazon.com/ls/docs/ja_jp/articles/amazon-lightsail-pushing-container-images
.PHONY: aws-push-image
aws-push-image: build-prod
	aws lightsail push-container-image --service-name ${LIGHTSAIL_SERVICE_NAME} --label webhook --image ${IMAGE_NAME}:${IMAGE_VERSION}

aws-deploy:
	aws lightsail create-container-service-deployment \
		--service-name ${LIGHTSAIL_SERVICE_NAME} \
		--containers '${LIGHTSAIL_CONTAINER_NAME}={image=${LIGHTSAIL_IMAGE_NAME_PREFIX}${LIGHTSAIL_IMAGE_VERSION},ports={80=HTTP}}' \
		--public-endpoint containerName=${LIGHTSAIL_CONTAINER_NAME},containerPort=80

test-request:
	curl \
		-XPOST \
		-H 'x-line-signature: ${TEST_REQUEST_SIGNATURE}' \
		--data '${TEST_REQUEST_DATA}' \
		localhost:8000/webhook
