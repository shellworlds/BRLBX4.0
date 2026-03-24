.PHONY: test build docker

test:
	cd ../.. && go test ./services/auth-rbac/...

build:
	cd ../.. && CGO_ENABLED=0 go build -trimpath -o /tmp/auth-rbac ./services/auth-rbac

docker:
	cd ../.. && docker build -f services/auth-rbac/Dockerfile -t auth-rbac:local .
