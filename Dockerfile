FROM golang:1.21.1

WORKDIR /detect-copies

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

CMD go run ./cmd/detect-copies -address $GROUP_ADDRESS -port $GROUP_PORT