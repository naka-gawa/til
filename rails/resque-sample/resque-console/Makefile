IMAGE_NAME?=resque-console
BUILD_TAG?=dev

.PHONY: build run
.ONSHELL:
SHELL=/bin/bash

build:
	docker build --rm -f Dockerfile -t ${IMAGE_NAME}:${BUILD_TAG} .

run:
	docker run -ti --link redis:redis --name ${IMAGE_NAME} ${IMAGE_NAME}:${BUILD_TAG}

