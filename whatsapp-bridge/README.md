# WhatsApp Bridge

A WhatsApp Web API bridge that allows you to interact with WhatsApp Web programmatically.

## Features

- WhatsApp Web client integration
- REST API for sending and receiving messages
- QR code authentication
- Persistent session storage
- Message history
- Real-time updates

## Prerequisites

- Go 1.22 or later
- Docker (optional)
- Railway account (for deployment)

## Local Development

1. Clone the repository:
```bash
git clone https://github.com/networkerman/whatsapp-web-ui.git
cd whatsapp-web-ui/whatsapp-bridge
```

2. Install dependencies:
```bash
go mod download
```

3. Run the server:
```bash
go run main.go
```

## Docker Development

1. Build the Docker image:
```bash
docker build -t whatsapp-bridge .
```

2. Run the container:
```bash
docker run -p 8080:8080 -v $(pwd)/store:/data/store whatsapp-bridge
```

## Deployment on Railway

1. Fork or clone the repository:
```bash
git clone https://github.com/networkerman/whatsapp-web-ui.git
```

2. Connect your GitHub repository to Railway:
   - Go to [Railway Dashboard](https://railway.app/dashboard)
   - Click "New Project"
   - Select "Deploy from GitHub repo"
   - Choose `whatsapp-web-ui` repository
   - Select the `api-server` branch

3. Configure environment variables in Railway:
   - `PORT`: 8080
   - `STORE_PATH`: /data/store
   - `STORE_PERMISSIONS`: 755

4. Add persistent storage:
   - Go to the "Volumes" tab
   - Create a new volume
   - Mount path: `/data/store`

The application will automatically deploy when you push changes to the `api-server` branch.

## API Endpoints

- `GET /api/status` - Check connection status
- `GET /api/qr` - Get QR code for authentication
- `GET /api/chats` - List all chats
- `GET /api/messages/{chatID}` - Get messages from a specific chat
- `POST /api/send` - Send a message

## License

MIT License 