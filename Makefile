VERSION=$(shell git describe --tags --always)
.PHONY: build
# build
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...
.PHONY: docker
docker:
	 docker build -t ccheers/go-blog:latest .
.PHONY: docker-amd64
docker-amd64:
	docker buildx build --platform linux/amd64 -t ccheers/go-blog:latest . --push
