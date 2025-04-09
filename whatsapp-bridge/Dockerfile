FROM golang:1.22 as builder

# Force cache invalidation with a timestamp
RUN echo "Build timestamp: $(date)" > /build_timestamp

WORKDIR /app

# Install required system dependencies for static compilation
RUN apt-get update && apt-get install -y \
    gcc \
    musl-tools \
    sqlite3 \
    libsqlite3-dev \
    libssl-dev \
    && rm -rf /var/lib/apt/lists/*

# Set up for static compilation with musl
RUN ln -s /usr/bin/musl-gcc /usr/local/bin/musl-gcc
ENV CC=musl-gcc

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

# Build a completely static binary with no dependencies
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=musl-gcc go build -ldflags="-linkmode external -extldflags '-static'" -tags "netgo,osusergo,static,sqlite_omit_load_extension" -o whatsapp-bridge cmd/server/main.go

# Verify it's truly static
RUN file whatsapp-bridge | grep "statically linked"
RUN ldd whatsapp-bridge || echo "No dynamic dependencies (this is good!)"

# Print binary info
RUN ls -l whatsapp-bridge
RUN sha256sum whatsapp-bridge

# Use scratch - completely empty base image
FROM scratch

# Copy TLS certificates for HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Create minimal application directory structure
WORKDIR /app

# Copy the static binary from builder
COPY --from=builder /app/whatsapp-bridge /app/whatsapp-bridge

# Create data directory in advance
WORKDIR /data

# No need to create users - scratch doesn't have users
# Set environment variables
ENV STORE_PATH=/data \
    STORE_PERMISSIONS=777 \
    TZ=UTC \
    PORT=8080 \
    DEBUG=true

# Expose the application port
EXPOSE 8080

# Run the binary directly - no entrypoint script needed in scratch
CMD ["/app/whatsapp-bridge"]

# No health check - scratch doesn't have curl