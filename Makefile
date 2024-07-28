# Image and tag can be overridden via environment variables.
DOCKER_USERNAME ?= synthao
IMAGE_NAME ?= imgproc
TAG ?= latest

APP_PORT?=7070

# Name of the Docker image.
IMAGE := ${DOCKER_USERNAME}/${IMAGE_NAME}:${TAG}

.PHONY: all gen

all: app-docker-build app-docker-push app-docker-run

gen:
	@protoc -I proto proto/imgproc/imgproc.proto --go_out=./gen/go/ --go_opt=paths=source_relative --go-grpc_out=./gen/go/ --go-grpc_opt=paths=source_relative

docker-build-push: app-docker-build app-docker-push

app-docker-build:
	@echo "Building Docker image ${IMAGE}"
	@docker build -t ${IMAGE} .

app-docker-push:
	@echo "Pushing Docker image ${IMAGE}"
	@docker push ${IMAGE}

app-docker-run:
	@echo "Running Docker image ${IMAGE_NAME}"
	@docker run -d --name $(IMAGE_NAME) -p $(APP_PORT):$(APP_PORT) --env-file .env $(IMAGE)