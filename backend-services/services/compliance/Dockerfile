# syntax=docker/dockerfile:1
FROM golang:1.22-alpine AS build
WORKDIR /src
RUN apk add --no-cache git ca-certificates
COPY go.mod go.sum ./
RUN go mod download
COPY pkg ./pkg
COPY services/compliance ./services/compliance
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /compliance ./services/compliance

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=build /src/services/compliance/migrations ./migrations
COPY --from=build /compliance /compliance
ENV SERVICE_ROOT=/app
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/compliance"]
