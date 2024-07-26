# Use a smaller base image
FROM golang:1.22-alpine AS build-dev

# Set the working directory
WORKDIR /app

# Set the env for the server
ENV GOOS="linux"
ENV CGO_ENABLED=0
ENV GO111MODULE="on"

# Copy go.mod and go.sum and install dependencies
RUN apk update \
  && apk add --no-cache \
  ca-certificates \
  curl \
  tzdata \
  git \
  && update-ca-certificates

# Install tools and dependencies
RUN go install github.com/air-verse/air@latest \
  && go install github.com/go-delve/delve/cmd/dlv@latest \
  && go install github.com/vektra/mockery/v2@latest

# Expose ports
EXPOSE ${HTTP_PORT}
EXPOSE 2345

# Use ENTRYPOINT to start the application
ENTRYPOINT ["air", "-c", ".air.toml"]