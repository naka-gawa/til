IMAGE_NAME?=redis
BUILD_TAG?=dev

.PHONY: build run
.ONSHELL:
SHELL=/bin/bash

build:
	docker build --rm -f Dockerfile -t ${IMAGE_NAME}:${BUILD_TAG} .

run:
	docker run --rm -d --name ${IMAGE_NAME} ${IMAGE_NAME}:${BUILD_TAG}

