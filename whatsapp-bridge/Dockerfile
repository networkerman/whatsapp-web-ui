FROM golang:1.22 as builder

# Force cache invalidation with a timestamp
RUN echo "Build timestamp: $(date)" > /build_timestamp

WORKDIR /app

# Install required system dependencies
RUN apt-get update && apt-get install -y \
    gcc \
    sqlite3 \
    libsqlite3-dev \
    libssl-dev \
    && rm -rf /var/lib/apt/lists/*

# Copy Go modules first for better caching
COPY go.mod go.sum ./
RUN go mod download

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
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o whatsapp-bridge cmd/server/main.go

# Use Debian instead of Alpine for better compatibility with CGO
FROM debian:bookworm-slim

# Force cache invalidation on runtime stage too
COPY --from=builder /build_timestamp /build_timestamp

# Install required runtime dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    sqlite3 \
    libsqlite3-dev \
    tzdata \
    openssl \
    libssl-dev \
    curl \
    procps \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/whatsapp-bridge .

# Set up the storage directory with proper permissions
RUN useradd -m -u 1000 appuser && \
    mkdir -p /data/store && \
    chown -R appuser:appuser /data/store /app && \
    chmod 777 /data/store

# Set environment variables
ENV STORE_PATH=/data/store \
    STORE_PERMISSIONS=755 \
    TZ=UTC

# Verify the binary has all required libraries
RUN chmod +x /app/whatsapp-bridge && \
    echo "Library dependencies:" && \
    ldd /app/whatsapp-bridge

# Copy the entrypoint script
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

# Switch to non-root user
USER appuser

# Expose the application port
EXPOSE 8080

# Set the entrypoint script as the container entrypoint
ENTRYPOINT ["/app/entrypoint.sh"]

# Define the command to run
CMD ["/app/whatsapp-bridge"]

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/api/status || exit 1