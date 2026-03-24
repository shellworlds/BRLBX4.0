# syntax=docker/dockerfile:1
FROM golang:1.22 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /out/ml-predictor ./services/ml-predictor
  
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=build /out/ml-predictor /app/ml-predictor
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/app/ml-predictor"]
