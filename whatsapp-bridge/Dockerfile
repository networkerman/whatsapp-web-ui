FROM golang:1.24-alpine

WORKDIR /app

# Install git and build dependencies
RUN apk add --no-cache git gcc musl-dev

# Create store directory with proper permissions
RUN mkdir -p /data/store && \
    chown -R nobody:nobody /data/store && \
    chmod -R 755 /data/store

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go mod tidy
RUN go build -o main .

EXPOSE 8080

# Environment variables for store configuration
ENV STORE_PATH=/data/store
ENV STORE_PERMISSIONS=755

# Create a non-root user
RUN adduser -D -g '' appuser
USER appuser

# Define volume for persistent storage
VOLUME ["/data/store"]

CMD ["./main"] 