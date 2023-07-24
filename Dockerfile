FROM golang:alpine
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

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
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.12
RUN go install github.com/vektra/mockery/v2@latest
RUN go install github.com/google/wire/cmd/wire@latest

EXPOSE 7070

# Run the air command in the directory where our code will live
ENTRYPOINT ["air", "-c", ".air.toml"]