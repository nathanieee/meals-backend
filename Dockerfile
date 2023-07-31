FROM golang:1.20-alpine
WORKDIR /app

ENV GO111MODULE="on"
ENV GOOS="linux"
ENV CGO_ENABLED=0
ENV GOFLAGS="-buildvcs=false"

RUN apk update \
  && apk add --no-cache \
  ca-certificates \
  curl \
  tzdata \
  git \
  && update-ca-certificates

RUN go install github.com/cosmtrek/air@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install github.com/vektra/mockery/v2@latest

COPY go.mod go.sum ./
RUN go mod download

EXPOSE 7070
EXPOSE 2345

ENTRYPOINT ["air", "-c", ".air.toml"]