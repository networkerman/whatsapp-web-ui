# WhatsApp Web Interface

A modern web interface for WhatsApp MCP, allowing you to interact with your WhatsApp messages through a beautiful web UI.

## Features

- View and manage your WhatsApp chats
- Send and receive messages
- Search through your message history
- Real-time message updates
- AI-powered message suggestions using Claude

## Prerequisites

- Node.js 18 or later
- npm or yarn
- WhatsApp MCP server running locally
- Anthropic Claude API key

## Setup

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd whatsapp-web-ui
   ```

2. Install dependencies for the frontend:
   ```bash
   cd frontend
   npm install
   ```

3. Install dependencies for the API server:
   ```bash
   cd ../api-server
   npm install
   ```

4. Configure environment variables:
   - Copy `.env.example` to `.env` in the `api-server` directory:
     ```bash
     cp .env.example .env
     ```
   - Update the values in `.env` with your configuration:
     - `ANTHROPIC_API_KEY`: Your Anthropic Claude API key
     - `MCP_SERVER_PATH`: Path to your WhatsApp MCP server

   ⚠️ **Security Note**: Never commit your `.env` file or any files containing API keys to version control. The `.env` file is already included in `.gitignore` to prevent accidental commits.

## Running the Application

1. Start the API server:
   ```bash
   cd api-server
   npm run dev
   ```

2. In a new terminal, start the frontend:
   ```bash
   cd frontend
   npm run dev
   ```

3. Open your browser and navigate to `http://localhost:3000`

## Development

- Frontend: Next.js with TypeScript and Tailwind CSS
- Backend: Express.js with TypeScript
- Real-time updates: Socket.IO
- AI Integration: Anthropic Claude API

## Security Best Practices

1. **API Keys**: Never commit API keys or sensitive credentials to version control
2. **Environment Variables**: Use `.env` files for local development and secure environment variable management for production
3. **Production Deployment**: Use your hosting platform's secure environment variable management (e.g., Vercel, Heroku, etc.)
4. **API Key Rotation**: Regularly rotate your API keys and update them in your environment variables

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details. 