# WhatsApp Bridge

A Node.js application that provides a WebSocket interface to interact with WhatsApp Web using Puppeteer.

## Features

- Automated WhatsApp Web interaction
- WebSocket server for sending and receiving messages
- QR code scanning for authentication
- Real-time status monitoring

## Prerequisites

- Node.js 18 or higher
- npm or yarn
- Chrome/Chromium browser

## Installation

1. Clone the repository
2. Install dependencies:
```bash
npm install
```

## Usage

1. Build the project:
```bash
npm run build
```

2. Start the server:
```bash
npm start
```

3. Open Chrome and scan the QR code that appears in the browser window

4. Connect to the WebSocket server (default port: 8080) and send messages using the following format:
```json
{
    "type": "send_message",
    "to": "1234567890",
    "message": "Hello, World!"
}
```

## WebSocket API

### Message Types

1. `send_message`: Send a message to a specific number
   ```json
   {
       "type": "send_message",
       "to": "1234567890",
       "message": "Hello, World!"
   }
   ```

2. `get_status`: Check the authentication status
   ```json
   {
       "type": "get_status"
   }
   ```

### Response Format

Success response:
```json
{
    "status": "authenticated"
}
```

Error response:
```json
{
    "error": "Error message here"
}
```

## Development

The project is written in TypeScript and uses:
- Puppeteer for browser automation
- ws for WebSocket server
- TypeScript for type safety

## License

MIT 