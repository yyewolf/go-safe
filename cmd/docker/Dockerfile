FROM golang:1.20-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/go-safe ./cmd/docker
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/go-safe-cli ./cmd/cli

FROM alpine
COPY --from=builder /app/go-safe /go-safe
COPY --from=builder /app/go-safe-cli /go-safe-cli
ENTRYPOINT ["/go-safe"]
