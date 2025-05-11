FROM golang:1.24.3-alpine AS builder

WORKDIR /app

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -o messaging-app ./cmd/server

# Use minimal image for final stage
FROM alpine:3.18

WORKDIR /app

# Copy built application
COPY --from=builder /app/messaging-app .

# Create app user
RUN adduser -D -g '' appuser
USER appuser

# Run
CMD ["./messaging-app"]
