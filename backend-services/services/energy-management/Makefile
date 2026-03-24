.PHONY: test build docker

test:
	cd ../.. && go test ./services/energy-management/...

build:
	cd ../.. && CGO_ENABLED=0 go build -trimpath -o /tmp/energy-management ./services/energy-management

docker:
	cd ../.. && docker build -f services/energy-management/Dockerfile -t energy-management:local .
