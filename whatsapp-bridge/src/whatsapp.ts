import puppeteer from 'puppeteer';
import { WebSocketServer, WebSocket } from 'ws';
import { EventEmitter } from 'events';

class WhatsAppClient extends EventEmitter {
    private browser: puppeteer.Browser | null = null;
    private page: puppeteer.Page | null = null;
    private wsServer: WebSocketServer | null = null;
    private isAuthenticated = false;

    constructor(private port: number = 8080) {
        super();
    }

    async start() {
        try {
            // Launch browser with specific options
            this.browser = await puppeteer.launch({
                headless: false,
                args: [
                    '--no-sandbox',
                    '--disable-setuid-sandbox',
                    '--disable-dev-shm-usage',
                    '--disable-accelerated-2d-canvas',
                    '--disable-gpu',
                    '--window-size=1920,1080'
                ]
            });

            // Create new page
            this.page = await this.browser.newPage();
            
            // Set viewport
            await this.page.setViewport({ width: 1920, height: 1080 });

            // Set user agent
            await this.page.setUserAgent('Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36');

            // Navigate to WhatsApp Web
            console.log('Navigating to WhatsApp Web...');
            await this.page.goto('https://web.whatsapp.com', {
                waitUntil: 'networkidle0',
                timeout: 60000
            });

            // Wait for QR code to be visible
            console.log('Waiting for QR code...');
            try {
                await this.page.waitForSelector('canvas', { visible: true, timeout: 60000 });
                console.log('QR code found!');
            } catch (error) {
                console.error('Error waiting for QR code:', error);
                throw error;
            }

            // Set up WebSocket server
            this.wsServer = new WebSocketServer({ port: this.port });
            console.log(`WebSocket server started on port ${this.port}`);

            this.wsServer.on('connection', (ws) => {
                console.log('Client connected to WebSocket');
                this.setupWebSocketHandlers(ws);
            });

            // Start monitoring for authentication
            this.monitorAuthentication();

        } catch (error) {
            console.error('Error starting WhatsApp client:', error);
            throw error;
        }
    }

    private async monitorAuthentication() {
        if (!this.page) return;

        try {
            // Wait for the main chat list to appear (indicating successful authentication)
            await this.page.waitForSelector('div[data-testid="chat-list"]', { timeout: 60000 });
            console.log('WhatsApp Web authenticated successfully!');
            this.isAuthenticated = true;
            this.emit('authenticated');
        } catch (error) {
            console.error('Error monitoring authentication:', error);
        }
    }

    private setupWebSocketHandlers(ws: WebSocket) {
        ws.on('message', async (message: Buffer) => {
            try {
                const data = JSON.parse(message.toString());
                console.log('Received message:', data);

                if (!this.isAuthenticated) {
                    ws.send(JSON.stringify({
                        error: 'WhatsApp not authenticated yet. Please scan the QR code first.'
                    }));
                    return;
                }

                // Handle different message types
                switch (data.type) {
                    case 'send_message':
                        await this.sendMessage(data.to, data.message);
                        break;
                    case 'get_status':
                        ws.send(JSON.stringify({
                            status: this.isAuthenticated ? 'authenticated' : 'waiting_for_qr'
                        }));
                        break;
                    default:
                        ws.send(JSON.stringify({
                            error: 'Unknown message type'
                        }));
                }
            } catch (error) {
                console.error('Error handling WebSocket message:', error);
                ws.send(JSON.stringify({
                    error: 'Internal server error'
                }));
            }
        });

        ws.on('close', () => {
            console.log('Client disconnected from WebSocket');
        });
    }

    private async sendMessage(to: string, message: string) {
        if (!this.page) return;

        try {
            // Navigate to the chat
            await this.page.goto(`https://web.whatsapp.com/send?phone=${to}`, {
                waitUntil: 'networkidle0',
                timeout: 30000
            });

            // Wait for the message input to be ready
            await this.page.waitForSelector('div[data-testid="conversation-compose-box-input"]', {
                timeout: 30000
            });

            // Type and send the message
            await this.page.type('div[data-testid="conversation-compose-box-input"]', message);
            await this.page.keyboard.press('Enter');

            // Wait for the message to be sent
            await this.page.waitForTimeout(2000);
        } catch (error) {
            console.error('Error sending message:', error);
            throw error;
        }
    }

    async stop() {
        if (this.wsServer) {
            this.wsServer.close();
            this.wsServer = null;
        }

        if (this.browser) {
            await this.browser.close();
            this.browser = null;
        }

        this.isAuthenticated = false;
    }
}

export default WhatsAppClient; 