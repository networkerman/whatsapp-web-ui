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

### 4. Frontend Setup

The frontend is a Next.js application that provides the web interface:

```bash
cd whatsapp-web-ui-frontend
npm install
npm run dev
```

## Environment Variables

### Frontend (Netlify)
```env
NEXT_PUBLIC_API_URL=https://your-railway-app.up.railway.app
```

### API Server (Railway)
```env
PORT=3001
FRONTEND_URL=https://messageai.netlify.app
ANTHROPIC_API_KEY=your_claude_api_key_here
MCP_SERVER_PATH=https://your-railway-app.up.railway.app
```

## Deployment

The project is set up for automatic deployment:

1. **GitHub Repository**: https://github.com/networkerman/whatsapp-web-ui.git
2. **Railway**: Automatically deploys the API server and WhatsApp bridge
3. **Netlify**: Hosts the frontend at https://messageai.netlify.app

### Environment Variables Setup

1. **In Railway**:
   - Go to your Railway project dashboard
   - Navigate to the "Variables" tab
   - Add the required environment variables:
     ```
     PORT=3001
     FRONTEND_URL=https://messageai.netlify.app
     ANTHROPIC_API_KEY=your_claude_api_key
     MCP_SERVER_PATH=https://your-railway-app.up.railway.app
     ```

2. **In Netlify**:
   - Go to your Netlify dashboard
   - Navigate to Site settings > Environment variables
   - Add:
     ```
     NEXT_PUBLIC_API_URL=https://your-railway-app.up.railway.app
     ```

## Troubleshooting

### 1. No Chats Showing
- Ensure the WhatsApp Bridge is running and connected
- Check the MCP server logs for any errors
- Verify the API server can connect to the MCP server
- Check the browser console for any errors
- Verify all environment variables are set correctly

### 2. Cannot Send Messages
- Check if the WhatsApp Bridge is still connected
- Verify the API server's MCP_SERVER_PATH is correct
- Check the browser console for any errors

### 3. Connection Issues
- Ensure all services are running
- Check environment variables are set correctly
- Verify network connectivity between components
- Check Railway and Netlify deployment logs

### 4. WhatsApp Bridge Issues
- If the QR code doesn't appear, try restarting the bridge
- If WhatsApp shows "Device Limit Reached", remove an existing device from WhatsApp settings
- If messages get out of sync, delete the database files and re-authenticate:
  ```bash
  rm whatsapp-bridge/store/messages.db
  rm whatsapp-bridge/store/whatsapp.db
  ```

## Security Considerations

1. The WhatsApp Bridge and MCP server should run on a secure, private network
2. Never expose the WhatsApp Bridge or MCP server directly to the internet
3. Use HTTPS for all external communications
4. Keep your API keys and credentials secure
5. Monitor your Railway and Netlify logs for any suspicious activity

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
