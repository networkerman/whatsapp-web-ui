FROM golang:1.23 as builder

WORKDIR /app

# Copy the entire source code first
COPY whatsapp-bridge/ .

# Initialize the module and download dependencies
RUN go mod init whatsapp-bridge && \
    go mod tidy && \
    go mod download

# Build the application
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o whatsapp-bridge

FROM alpine:3.19

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=builder /app/whatsapp-bridge .

RUN adduser -D -u 1000 appuser && \
    mkdir -p /data/store && \
    chown -R appuser:appuser /data/store && \
    chmod 755 /data/store

ENV STORE_PATH=/data/store
ENV STORE_PERMISSIONS=755

USER appuser
EXPOSE 8080

CMD ["./whatsapp-bridge"] 