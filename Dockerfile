FROM golang:1.24 as builder

WORKDIR /app

# Install build dependencies
RUN apt-get update && apt-get install -y \
    gcc \
    sqlite3 \
    libsqlite3-dev \
    libssl-dev \
    && rm -rf /var/lib/apt/lists/*

# Create dummy store package to satisfy imports
RUN mkdir -p /app/internal/store
RUN echo 'package store' > /app/internal/store/store.go

# Copy the entire source code first
COPY whatsapp-bridge/ .

# Download dependencies and build
RUN go mod download && go mod tidy

# Print directory structure for debugging
RUN ls -la
RUN find . -name '*.go' | grep -v vendor

# Build the application with SQLite support
RUN CGO_ENABLED=1 GOOS=linux go build -o whatsapp-bridge cmd/server/main.go

FROM debian:bookworm-slim

# Install runtime dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    sqlite3 \
    openssl \
    curl \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY --from=builder /app/whatsapp-bridge .

RUN useradd -m -u 1000 appuser && \
    mkdir -p /data/store && \
    chown -R appuser:appuser /data/store && \
    chmod 777 /data/store

# Set environment variables
ENV STORE_PATH=/data/store \
    STORE_PERMISSIONS=777 \
    PORT=8080 \
    DEBUG=true

# Test binary executable
RUN chmod +x /app/whatsapp-bridge && \
    ls -l /app/whatsapp-bridge && \
    ldd /app/whatsapp-bridge

USER appuser
EXPOSE 8080

# Define the command to run
CMD ["/app/whatsapp-bridge"]