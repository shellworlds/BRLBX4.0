# syntax=docker/dockerfile:1
FROM golang:1.22-alpine AS build
WORKDIR /src
RUN apk add --no-cache git ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY pkg ./pkg
COPY services/vendor-ecosystem ./services/vendor-ecosystem
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /vendor ./services/vendor-ecosystem

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=build /src/services/vendor-ecosystem/migrations ./migrations
COPY --from=build /vendor /vendor
ENV SERVICE_ROOT=/app
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/vendor"]
