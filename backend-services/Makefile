SHELL := /bin/bash
export GO111MODULE ?= on

.PHONY: test cover docker-all

test:
	go test ./...

cover:
	go test -cover ./...

docker-all:
	@test -n "$(TAG)" || (echo "set TAG=..." && exit 1)
	docker build -f services/hello-world/Dockerfile -t hello-world:$(TAG) .
	docker build -f services/energy-management/Dockerfile -t energy-management:$(TAG) .
	docker build -f services/vendor-ecosystem/Dockerfile -t vendor-ecosystem:$(TAG) .
	docker build -f services/iot-ingestion/Dockerfile -t iot-ingestion:$(TAG) .
	docker build -f services/auth-rbac/Dockerfile -t auth-rbac:$(TAG) .
