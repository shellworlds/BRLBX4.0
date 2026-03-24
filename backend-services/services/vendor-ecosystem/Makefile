.PHONY: test build docker

test:
	cd ../.. && go test ./services/vendor-ecosystem/...

build:
	cd ../.. && CGO_ENABLED=0 go build -trimpath -o /tmp/vendor-ecosystem ./services/vendor-ecosystem

docker:
	cd ../.. && docker build -f services/vendor-ecosystem/Dockerfile -t vendor-ecosystem:local .
