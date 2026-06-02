# Build stage
FROM golang:1.26.3-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o user-management ./cmd/main.go

# Run stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/user-management .
COPY --from=builder /app/.env.example .env 2>/dev/null || true

# Expose port
EXPOSE 8081

# Run
CMD ["./user-management"]