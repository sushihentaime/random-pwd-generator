DOCKER_IMAGE_NAME = pwdapp
DOCKER_TAG = latest
DOCKERFILE_PATH = .
CONTAINER_NAME = pwdapp_container
PORT = 4000

.PHONY: build
build:
	docker build -t $(DOCKER_IMAGE_NAME):$(DOCKER_TAG) $(DOCKERFILE_PATH)

.PHONY: run
run:
	docker run --rm -d --name $(CONTAINER_NAME) -p $(PORT):4000 $(DOCKER_IMAGE_NAME):$(DOCKER_TAG)

.PHONY: stop
stop:
	docker stop $(CONTAINER_NAME)

.PHONY: clean
clean:
	docker rmi $(DOCKER_IMAGE_NAME):$(DOCKER_TAG)
