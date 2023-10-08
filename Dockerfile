FROM golang:1.18

WORKDIR /detect-copies

COPY go.mod go.sum ./
COPY cmd ./cmd
COPY internal ./internal

RUN go mod download
ENTRYPOINT go run ./cmd/detect-copies -address $GROUP_ADDRESS -port $GROUP_PORT