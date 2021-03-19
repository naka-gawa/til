IMAGE_NAME?=resque-worker
BUILD_TAG?=dev

.PHONY: build run
.ONSHELL:
SHELL=/bin/bash

build:
	docker build --rm -f Dockerfile -t ${IMAGE_NAME}:${BUILD_TAG} .

run:
	docker run -d --link redis:redis --volumes-from resque-console --name="resque-worker1" ${IMAGE_NAME}:${BUILD_TAG}
	docker run -d --link redis:redis --volumes-from resque-console --name="resque-worker2" ${IMAGE_NAME}:${BUILD_TAG}

