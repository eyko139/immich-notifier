# syntax=docker/dockerfile:1

FROM golang:1.23.2

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./ui ./ui


# Build
RUN go build -o notifier ./cmd/web

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose

# Run
CMD ["./notifier"]
