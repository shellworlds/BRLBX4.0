.PHONY: test build docker

test:
	cd ../.. && go test ./services/ml-predictor/...

build:
	cd ../.. && CGO_ENABLED=0 go build -trimpath -o /tmp/ml-predictor ./services/ml-predictor

docker:
	cd ../.. && docker build -f services/ml-predictor/Dockerfile -t ml-predictor:local .
