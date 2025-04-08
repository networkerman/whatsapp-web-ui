FROM golang:1.24-alpine

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Create store directory with proper permissions
RUN mkdir -p /data/store && \
    chown -R nobody:nobody /data/store && \
    chmod 755 /data/store

# Build the application
RUN CGO_ENABLED=1 go build -o main .

# Set environment variables
ENV STORE_PATH=/data/store
ENV STORE_PERMISSIONS=755

# Create non-root user
RUN adduser -D -u 1000 appuser && \
    chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Define volume for persistent storage
VOLUME ["/data/store"]

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"] 