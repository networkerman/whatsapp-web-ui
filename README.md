# WhatsApp MCP with Web UI

This project combines a WhatsApp MCP (Model Context Protocol) server with a modern web interface. It allows you to interact with your WhatsApp messages through a beautiful web UI while maintaining privacy and security.

## Project Structure

```
whatsapp-mcp/
├── whatsapp-bridge/      # Go application that connects to WhatsApp
├── whatsapp-mcp-server/  # Python MCP server that manages WhatsApp data
├── whatsapp-web-ui/      # Next.js frontend application
│   ├── frontend/        # Frontend code
│   └── api-server/      # Go API server
└── whatsapp-web-ui-frontend/  # Deployed frontend (separate repo)
```

## Prerequisites

- Go 1.20 or later
- Python 3.8 or later
- Docker and Docker Compose
- Node.js 18 or later
- UV (Python package manager): `curl -LsSf https://astral.sh/uv/install.sh | sh`

## Setup Instructions

### 1. WhatsApp Bridge Setup

The WhatsApp Bridge is responsible for connecting to your WhatsApp account:

```bash
cd whatsapp-bridge
go run main.go
```

On first run, you'll see a QR code in the terminal. Scan this with your WhatsApp mobile app:
1. Open WhatsApp on your phone
2. Go to Settings → Linked Devices
3. Tap "Link a Device"
4. Scan the QR code shown in the terminal

The bridge will maintain this connection and store the session locally.

### 2. WhatsApp MCP Server Setup

The MCP server manages the WhatsApp connection and stores messages:

```bash
cd whatsapp-mcp-server
uv venv
source .venv/bin/activate  # On Windows: .venv\Scripts\activate
uv pip install -r requirements.txt
python main.py
```

### 3. API Server Setup

The API server communicates between the frontend and the MCP server:

```bash
cd whatsapp-web-ui/api-server
go run main.go
```

Set these environment variables:
```env
PORT=3001
FRONTEND_URL=http://localhost:3000
ANTHROPIC_API_KEY=your_claude_api_key_here
MCP_SERVER_PATH=/path/to/whatsapp-mcp-server
```

### 4. Frontend Setup

The frontend is a Next.js application that provides the web interface:

```bash
cd whatsapp-web-ui-frontend
npm install
npm run dev
```

Set these environment variables:
```env
NEXT_PUBLIC_API_URL=http://localhost:3001
```

## Deployment

### Frontend (Netlify)
1. Create a new repository on GitHub
2. Push the frontend code to the repository
3. Connect the repository to Netlify
4. Set environment variables in Netlify:
   - `NEXT_PUBLIC_API_URL`: Your API server URL

### API Server (Railway)
1. Push the API server code to GitHub
2. Connect to Railway
3. Set environment variables in Railway:
   - `PORT`: 3001
   - `FRONTEND_URL`: Your frontend URL
   - `ANTHROPIC_API_KEY`: Your Claude API key
   - `MCP_SERVER_PATH`: Path to your MCP server

### WhatsApp Bridge and MCP Server
These components should run locally or on a server you control, as they need direct access to your WhatsApp account.

## Environment Variables

### Frontend (.env.local)
```env
NEXT_PUBLIC_API_URL=http://localhost:3001
```

### API Server (.env)
```env
PORT=3001
FRONTEND_URL=http://localhost:3000
ANTHROPIC_API_KEY=your_claude_api_key_here
MCP_SERVER_PATH=/path/to/whatsapp-mcp-server
```

## Troubleshooting

1. **No Chats Showing**
   - Ensure the WhatsApp Bridge is running and connected
   - Check the MCP server logs for any errors
   - Verify the API server can connect to the MCP server

2. **Cannot Send Messages**
   - Check if the WhatsApp Bridge is still connected
   - Verify the API server's MCP_SERVER_PATH is correct
   - Check the browser console for any errors

3. **Connection Issues**
   - Ensure all services are running
   - Check environment variables are set correctly
   - Verify network connectivity between components

## Security Considerations

1. The WhatsApp Bridge and MCP server should run on a secure, private network
2. Never expose the WhatsApp Bridge or MCP server directly to the internet
3. Use HTTPS for all external communications
4. Keep your API keys and credentials secure

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
