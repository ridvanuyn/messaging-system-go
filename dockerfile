FROM golang:1.24.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o messaging-app ./cmd/server

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/messaging-app .

RUN adduser -D -g '' appuser
USER appuser

CMD ["./messaging-app"]
