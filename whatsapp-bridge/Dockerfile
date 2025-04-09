FROM golang:1.22 as builder

WORKDIR /app

# Install required system dependencies
RUN apt-get update && apt-get install -y \
    gcc \
    sqlite3 \
    libssl-dev \
    && rm -rf /var/lib/apt/lists/*

# Copy only go.mod and go.sum first for better caching
COPY go.mod go.sum ./

# Prepare directories for compilation
RUN mkdir -p /app/internal/store
RUN touch /app/internal/store/store.go
RUN echo 'package store' > /app/internal/store/store.go

# Copy the rest of the source code
COPY . .

# Verify module structure
RUN find . -type f -name '*.go' | xargs cat | grep -n 'import' || true

# Patch imports before building
RUN find . -type f -name '*.go' | xargs grep -l 'internal/store' || true

# Build the application with SQLite support and enable debugging output
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -tags "debug" -o whatsapp-bridge cmd/server/main.go

FROM alpine:3.19

# Install required runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    sqlite \
    tzdata \
    openssl \
    curl

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/whatsapp-bridge .

# Set up the storage directory with proper permissions
RUN adduser -D -u 1000 appuser && \
    mkdir -p /data/store && \
    chown -R appuser:appuser /data/store /app && \
    chmod 777 /data/store

# Set environment variables
ENV STORE_PATH=/data/store \
    STORE_PERMISSIONS=755 \
    TZ=UTC

# Switch to non-root user
USER appuser

# Expose the application port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/api/status || exit 1

CMD ["./whatsapp-bridge"]