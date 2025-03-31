import Anthropic from '@anthropic-ai/sdk';
import { spawn } from 'child_process';
import path from 'path';

export interface Chat {
  id: string;
  name: string;
  lastMessage: string;
  timestamp: string;
  unreadCount: number;
}

export interface Message {
  id: string;
  content: string;
  sender: string;
  timestamp: string;
  isOutgoing: boolean;
}

export class WhatsAppService {
  private anthropic: Anthropic;
  private mcpServerPath: string;
  private mcpProcess: any;

  constructor(apiKey: string, mcpServerPath: string) {
    this.anthropic = new Anthropic({ apiKey });
    this.mcpServerPath = mcpServerPath;
    this.mcpProcess = null;
  }

  async startMCPServer(): Promise<void> {
    return new Promise((resolve, reject) => {
      const mcpServerDir = path.dirname(this.mcpServerPath);
      const mcpServerFile = path.basename(this.mcpServerPath);

      this.mcpProcess = spawn('python3', [mcpServerFile], {
        cwd: mcpServerDir,
        stdio: ['pipe', 'pipe', 'pipe']
      });

      this.mcpProcess.stdout.on('data', (data: Buffer) => {
        console.log(`MCP Server: ${data.toString()}`);
        if (data.toString().includes('Server started')) {
          resolve();
        }
      });

      this.mcpProcess.stderr.on('data', (data: Buffer) => {
        console.error(`MCP Server Error: ${data.toString()}`);
        reject(new Error(data.toString()));
      });

      this.mcpProcess.on('close', (code: number) => {
        console.log(`MCP Server exited with code ${code}`);
      });
    });
  }

  async stopMCPServer(): Promise<void> {
    if (this.mcpProcess) {
      this.mcpProcess.kill();
      this.mcpProcess = null;
    }
  }

  async getChats(): Promise<Chat[]> {
    try {
      // Use Claude to get chat list through MCP server
      const response = await this.anthropic.messages.create({
        model: 'claude-3-sonnet-20240229',
        max_tokens: 1000,
        messages: [{
          role: 'user',
          content: 'List all my WhatsApp chats'
        }]
      });

      // Parse the response and convert to Chat interface
      // This is a placeholder - you'll need to implement the actual parsing
      return [];
    } catch (error) {
      console.error('Error getting chats:', error);
      throw error;
    }
  }

  async getMessages(chatId: string): Promise<Message[]> {
    try {
      // Use Claude to get messages through MCP server
      const response = await this.anthropic.messages.create({
        model: 'claude-3-sonnet-20240229',
        max_tokens: 1000,
        messages: [{
          role: 'user',
          content: `Get messages for chat ${chatId}`
        }]
      });

      // Parse the response and convert to Message interface
      // This is a placeholder - you'll need to implement the actual parsing
      return [];
    } catch (error) {
      console.error('Error getting messages:', error);
      throw error;
    }
  }

  async sendMessage(chatId: string, message: string): Promise<void> {
    try {
      // Use Claude to send message through MCP server
      await this.anthropic.messages.create({
        model: 'claude-3-sonnet-20240229',
        max_tokens: 1000,
        messages: [{
          role: 'user',
          content: `Send message "${message}" to chat ${chatId}`
        }]
      });
    } catch (error) {
      console.error('Error sending message:', error);
      throw error;
    }
  }

  async searchMessages(query: string): Promise<Message[]> {
    try {
      // Use Claude to search messages through MCP server
      const response = await this.anthropic.messages.create({
        model: 'claude-3-sonnet-20240229',
        max_tokens: 1000,
        messages: [{
          role: 'user',
          content: `Search messages for "${query}"`
        }]
      });

      // Parse the response and convert to Message interface
      // This is a placeholder - you'll need to implement the actual parsing
      return [];
    } catch (error) {
      console.error('Error searching messages:', error);
      throw error;
    }
  }
} 