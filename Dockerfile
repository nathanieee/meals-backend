# Use a smaller base image
FROM golang:1.21-alpine AS build-dev

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum and install dependencies
COPY go.mod go.sum ./
RUN apk update && apk add --no-cache ca-certificates curl tzdata git \
  && update-ca-certificates \
  && go mod download

# Install tools and dependencies
RUN go install github.com/cosmtrek/air@latest \
  && go install github.com/go-delve/delve/cmd/dlv@latest \
  && go install github.com/vektra/mockery/v2@latest

# Copy the rest of the source code
COPY . .

# Expose ports
EXPOSE ${HTTP_PORT}
EXPOSE 2345

# Use ENTRYPOINT to start the application
ENTRYPOINT ["air", "-c", ".air.toml"]
