FROM golang:1.24-alpine

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Create store directory with correct permissions
RUN mkdir -p /data/store && \
    chown -R nobody:nobody /data/store && \
    chmod 755 /data/store

# Set environment variables
ENV STORE_PATH=/data/store
ENV STORE_PERMISSIONS=755

# Build the application
RUN go build -o main .

# Create a non-root user
RUN adduser -D -g '' appuser && \
    chown -R appuser:appuser /app /data/store

# Switch to non-root user
USER appuser

# Expose the port
EXPOSE 8080

# Command to run the executable
CMD ["./main"] 