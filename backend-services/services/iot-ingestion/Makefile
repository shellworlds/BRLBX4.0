.PHONY: test build docker

test:
	cd ../.. && go test ./services/iot-ingestion/...

build:
	cd ../.. && CGO_ENABLED=0 go build -trimpath -o /tmp/iot-ingestion ./services/iot-ingestion

docker:
	cd ../.. && docker build -f services/iot-ingestion/Dockerfile -t iot-ingestion:local .
