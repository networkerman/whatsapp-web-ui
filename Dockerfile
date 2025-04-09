FROM golang:1.24 as builder

WORKDIR /app

# Install build dependencies
RUN apt-get update && apt-get install -y \
    gcc-x86-64-linux-gnu \
    libssl-dev \
    && rm -rf /var/lib/apt/lists/*

# Copy the entire source code first
COPY whatsapp-bridge/ .

# Download dependencies and build
RUN go mod tidy && \
    go mod download

# Build the application
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=x86_64-linux-gnu-gcc go build -o whatsapp-bridge

FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    openssl \
    curl

WORKDIR /app
COPY --from=builder /app/whatsapp-bridge .

RUN adduser -D -u 1000 appuser && \
    mkdir -p /data/store && \
    chown -R appuser:appuser /data/store && \
    chmod 777 /data/store

ENV STORE_PATH=/data/store
ENV STORE_PERMISSIONS=777

USER appuser
EXPOSE 8080

CMD ["./whatsapp-bridge"] 