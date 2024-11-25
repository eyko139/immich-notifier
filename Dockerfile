# syntax=docker/dockerfile:1

FROM golang:1.23.2 AS builder

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./ui ./ui

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o notifier ./cmd/web

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/notifier .
COPY --from=builder /app/ui ./ui

ARG APP_VERSION
ENV APP_VERSION=${APP_VERSION}

# Run
CMD ["./notifier"]
