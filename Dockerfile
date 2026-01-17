# syntax=docker/dockerfile:1

# Build the application source
FROM golang:latest AS build-stage

WORKDIR /app

COPY . .

RUN go mod download

RUN go test ./...

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/cc-webserver .

# Deploy the application binary into a lean image
FROM alpine:latest AS build-release-stage

WORKDIR /app

COPY --from=build-stage /app/cc-webserver .

ENTRYPOINT ["./cc-webserver", "-addr=0.0.0.0:8080"]
