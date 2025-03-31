import express from 'express';
import cors from 'cors';
import { createServer } from 'http';
import { Server } from 'socket.io';
import dotenv from 'dotenv';
import Anthropic from '@anthropic-ai/sdk';

dotenv.config();

const app = express();
const httpServer = createServer(app);
const io = new Server(httpServer, {
  cors: {
    origin: process.env.FRONTEND_URL || 'http://localhost:3000',
    methods: ['GET', 'POST']
  }
});

// Initialize Anthropic client
const anthropic = new Anthropic({
  apiKey: process.env.ANTHROPIC_API_KEY
});

app.use(cors());
app.use(express.json());

// Mock data for testing
const mockChats = [
  { id: '1', name: 'Family Group', lastMessage: 'See you tomorrow!', timestamp: Date.now() - 3600000 },
  { id: '2', name: 'Work Team', lastMessage: 'Meeting at 2 PM', timestamp: Date.now() - 7200000 },
  { id: '3', name: 'Friends', lastMessage: 'Party tonight?', timestamp: Date.now() - 86400000 },
];

const mockMessages: Record<string, any[]> = {
  '1': [
    { id: '1', content: 'Hi everyone!', timestamp: Date.now() - 3600000, sender: 'user' },
    { id: '2', content: 'Hello!', timestamp: Date.now() - 3500000, sender: 'bot' },
  ],
  '2': [
    { id: '3', content: 'Project update?', timestamp: Date.now() - 7200000, sender: 'user' },
    { id: '4', content: 'Everything is on track', timestamp: Date.now() - 7100000, sender: 'bot' },
  ],
  '3': [
    { id: '5', content: 'Anyone up for dinner?', timestamp: Date.now() - 86400000, sender: 'user' },
    { id: '6', content: 'Count me in!', timestamp: Date.now() - 86300000, sender: 'bot' },
  ],
};

// Routes
app.get('/api/chats', async (req, res) => {
  try {
    res.json(mockChats);
  } catch (error) {
    console.error('Error fetching chats:', error);
    res.status(500).json({ error: 'Failed to fetch chats' });
  }
});

app.get('/api/messages/:chatId', async (req, res) => {
  try {
    const { chatId } = req.params;
    const messages = mockMessages[chatId] || [];
    res.json(messages);
  } catch (error) {
    console.error('Error fetching messages:', error);
    res.status(500).json({ error: 'Failed to fetch messages' });
  }
});

app.post('/api/messages/:chatId', async (req, res) => {
  try {
    const { chatId } = req.params;
    const { message } = req.body;

    // Create a new message
    const newMessage = {
      id: Date.now().toString(),
      content: message,
      timestamp: Date.now(),
      sender: 'user' as const
    };

    // Add to mock messages
    if (!mockMessages[chatId]) {
      mockMessages[chatId] = [];
    }
    mockMessages[chatId].push(newMessage);

    // Update last message in chat
    const chat = mockChats.find(c => c.id === chatId);
    if (chat) {
      chat.lastMessage = message;
      chat.timestamp = Date.now();
    }

    // Emit the new message to connected clients
    io.emit(`chat:${chatId}`, newMessage);

    res.json(newMessage);
  } catch (error) {
    console.error('Error sending message:', error);
    res.status(500).json({ error: 'Failed to send message' });
  }
});

// WebSocket connection handling
io.on('connection', (socket) => {
  console.log('Client connected');

  socket.on('disconnect', () => {
    console.log('Client disconnected');
  });

  // Handle real-time message updates
  socket.on('subscribe_to_chat', (chatId: string) => {
    socket.join(`chat:${chatId}`);
    console.log(`Client subscribed to chat: ${chatId}`);
  });
});

const PORT = process.env.PORT || 3001;
httpServer.listen(PORT, () => {
  console.log(`API server running on port ${PORT}`);
}); 