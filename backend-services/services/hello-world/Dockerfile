# syntax=docker/dockerfile:1
FROM golang:1.22-alpine AS build
WORKDIR /src
RUN apk add --no-cache ca-certificates git
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /hello-world .

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /
COPY --from=build /hello-world /hello-world
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/hello-world"]
