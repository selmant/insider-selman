FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o mock_webhook ./cmd/mock_webhook.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/.env .
COPY --from=builder /app/mock_webhook .

EXPOSE 8080
CMD ["./mock_webhook"]