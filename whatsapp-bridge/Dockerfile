# Use the official Go image with Alpine
FROM golang:1.22-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev openssl-dev

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./

# Create dummy store package to satisfy imports
RUN mkdir -p /app/internal/store
RUN echo 'package store' > /app/internal/store/store.go

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /app/whatsapp-bridge cmd/server/main.go

# Use Alpine for runtime
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata sqlite openssl

# Create app directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/whatsapp-bridge /app/whatsapp-bridge

# Create data directory with proper permissions
RUN mkdir -p /data/store && chmod 777 /data/store

# Set environment variables
ENV STORE_PATH=/data/store \
    STORE_PERMISSIONS=777 \
    TZ=UTC \
    PORT=8080 \
    DEBUG=true

# Expose the port
EXPOSE 8080

# Define health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget -q --spider http://localhost:8080/api/status || exit 1

# Run the application
CMD ["/app/whatsapp-bridge"]