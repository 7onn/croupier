default: help

help:
	@echo 
	@echo "available commands"
	@echo "  - dev"
	@echo "  - push-latest"
	@echo "  - clean"
	@echo
	
GOLANG_DOCKER_IMAGE := golang:1.15.6
GO_BUILD_CONTAINER=docker run -i --rm \
	-v "$$(pwd):/ops" \
	-w /ops --user="$$(id -u):$$(id -g)" \
	-e GOPATH=/ops/build-go-path \
	-e GOCACHE=/ops/build-cache \
	$(GOLANG_DOCKER_IMAGE)

BIN=server
GO=CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go
GO_BUILD_COMMAND=$(GO) build -a --ldflags '-X main.VERSION=$(TAG) -w -extldflags "-static"' -tags netgo -o $(BIN) ./

TAG=0.1.0
IMAGE=devbytom/pk-$(BIN)

.PHONY: dev
dev:
	go mod download
	go mod vendor
	air .

.PHONY: build
build:
	$(GO_BUILD_COMMAND)

.PHONY: docker-build
docker-build:
	$(GO_BUILD_CONTAINER) \
	sh -c "$(GO_BUILD_COMMAND)"

.PHONY: image
image: docker-build
	docker build -t $(IMAGE):$(TAG) .

.PHONY: push
push: image
	docker push $(IMAGE):$(TAG)

.PHONY: push-latest
push-latest: image
	docker tag $(IMAGE):$(TAG) $(IMAGE):latest
	docker push $(IMAGE):latest

.PHONY: clean
clean:
	rm $(BIN)
