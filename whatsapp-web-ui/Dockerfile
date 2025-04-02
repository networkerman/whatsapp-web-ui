# Use the official Golang image to create a build artifact.
FROM golang:1.21-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY api-server/go.mod api-server/go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY api-server/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Use a minimal image for the final stage
FROM alpine:latest

# Install CA certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

# Expose port
EXPOSE 8080

# Set environment variables
ENV PORT=8080

# Run the binary
CMD ["./main"] 