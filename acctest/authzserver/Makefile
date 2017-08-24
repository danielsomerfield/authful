.PHONY: all build clean
all: push

STAMP=`date +'%Y-%m-%d_%H%M%S'`
APP_NAME=authzserver
REGISTRY_USERNAME=danielsomerfield

build:
	cp ../../authful-linux-386 .
	docker build --tag "${REGISTRY_USERNAME}/${APP_NAME}:latest" --tag "${REGISTRY_USERNAME}/${APP_NAME}:${STAMP}" .

push: build
	docker push ${REGISTRY_USERNAME}/${APP_NAME}

clean:
	docker rmi -f ${REGISTRY_USERNAME}/${APP_NAME}