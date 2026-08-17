SHELL := /bin/sh
IMAGE_GO := golang:1.23.12
BIN := bin/transitmanifest

.PHONY: local docker-test docker-vet build clean

local:
	CGO_ENABLED=0 go test ./...
	CGO_ENABLED=0 go vet ./...
	CGO_ENABLED=0 go build -trimpath -o $(BIN) ./cmd/transitmanifest

docker-test:
	docker run --rm -v "$(CURDIR):/src" -w /src $(IMAGE_GO) go test ./...

docker-vet:
	docker run --rm -v "$(CURDIR):/src" -w /src -e CGO_ENABLED=0 $(IMAGE_GO) go vet ./...

build:
	./scripts/build.sh

clean:
	rm -rf bin

