# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install dependencies required for CGO/SQLite
RUN apk add --no-cache \
    gcc \
    g++ \
    musl-dev \
    sqlite \
    sqlite-dev \
    pkgconfig

# Copy the whatsapp-bridge directory
COPY whatsapp-bridge/ .

# Print directory structure for debugging
RUN ls -la && \
    find . -name '*.go' | sort

# Download dependencies
RUN go mod download && go mod tidy

# Build with CGO enabled for SQLite support
RUN CGO_ENABLED=1 GOOS=linux go build -o whatsapp-bridge cmd/server/main.go

# Test if binary is executable and check library dependencies
RUN chmod +x /app/whatsapp-bridge && \
    ls -l /app/whatsapp-bridge && \
    ldd /app/whatsapp-bridge

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    sqlite \
    sqlite-libs \
    curl

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/whatsapp-bridge .

# Create data directory with appropriate permissions
RUN mkdir -p /data/store && \
    chmod -R 777 /data/store

# Create startup script to ensure data directory exists
RUN printf '#!/bin/bash\n\
echo "Starting WhatsApp Bridge..."\n\
echo "Data directory: $STORE_PATH"\n\
mkdir -p "$STORE_PATH"\n\
chmod -R "$STORE_PERMISSIONS" "$STORE_PATH"\n\
echo "Directory contents:"\n\
ls -la "$STORE_PATH"\n\
exec /app/whatsapp-bridge\n' > /app/start.sh && \
    chmod +x /app/start.sh

# Expose port
EXPOSE 8080

# Define environment variables
ENV STORE_PATH=/data/store \
    STORE_PERMISSIONS=777 \
    PORT=8080 \
    DEBUG=true

# Define the command to run
CMD ["/app/start.sh"]