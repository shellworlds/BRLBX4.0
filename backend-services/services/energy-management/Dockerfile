# syntax=docker/dockerfile:1
# docker build -f backend-services/services/energy-management/Dockerfile -t energy-management:dev backend-services
FROM golang:1.22-alpine AS build
WORKDIR /src
RUN apk add --no-cache git ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY pkg ./pkg
COPY services/energy-management ./services/energy-management
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /energy ./services/energy-management

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=build /src/services/energy-management/migrations ./migrations
COPY --from=build /energy /energy
ENV SERVICE_ROOT=/app
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/energy"]
