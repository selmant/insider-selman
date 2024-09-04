FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o message_sender ./cmd/message_sender.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/.env .
COPY --from=builder /app/message_sender .

EXPOSE 8081
CMD ["./message_sender"]