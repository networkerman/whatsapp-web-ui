const API_URL = process.env.NEXT_PUBLIC_API_URL;

export interface Chat {
  id: string;
  name: string;
  lastMessage?: string;
  timestamp: number;
}

export interface Message {
  id: string;
  content: string;
  timestamp: number;
  sender: 'user' | 'bot';
}

export async function getChats(): Promise<Chat[]> {
  const response = await fetch(`${API_URL}/api/chats`);
  if (!response.ok) {
    throw new Error('Failed to fetch chats');
  }
  return response.json();
}

export async function getMessages(chatId: string): Promise<Message[]> {
  const response = await fetch(`${API_URL}/api/messages/${chatId}`);
  if (!response.ok) {
    throw new Error('Failed to fetch messages');
  }
  return response.json();
}

export async function sendMessage(chatId: string, message: string): Promise<Message> {
  const response = await fetch(`${API_URL}/api/messages/${chatId}`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ message }),
  });
  if (!response.ok) {
    throw new Error('Failed to send message');
  }
  return response.json();
} 