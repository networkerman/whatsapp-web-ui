const puppeteer = require('puppeteer');
const WebSocket = require('ws');
const qrcode = require('qrcode-terminal');

class WhatsAppWeb {
    constructor() {
        this.browser = null;
        this.page = null;
        this.wss = null;
        this.clients = new Set();
    }

    async start() {
        // Start WebSocket server
        this.wss = new WebSocket.Server({ port: 8081 });
        console.log('WebSocket server started on port 8081');

        this.wss.on('connection', (ws) => {
            this.clients.add(ws);
            console.log('Client connected');

            ws.on('message', async (message) => {
                try {
                    const data = JSON.parse(message);
                    if (data.type === 'message') {
                        await this.sendMessage(data.data.to, data.data.content);
                    }
                } catch (error) {
                    console.error('Error handling message:', error);
                }
            });

            ws.on('close', () => {
                this.clients.delete(ws);
                console.log('Client disconnected');
            });
        });

        // Launch browser
        this.browser = await puppeteer.launch({
            headless: "new",
            args: ['--no-sandbox', '--disable-setuid-sandbox']
        });

        // Create new page
        this.page = await this.browser.newPage();

        // Navigate to WhatsApp Web
        await this.page.goto('https://web.whatsapp.com');
        console.log('Navigated to WhatsApp Web');

        // Wait for QR code
        const qrCodeSelector = '[data-testid="qrcode"]';
        await this.page.waitForSelector(qrCodeSelector);
        console.log('QR code found');

        // Get QR code data
        const qrCodeData = await this.page.evaluate((selector) => {
            const element = document.querySelector(selector);
            return element ? element.getAttribute('data-ref') : null;
        }, qrCodeSelector);

        if (qrCodeData) {
            // Display QR code in terminal
            qrcode.generate(qrCodeData, { small: true });
            console.log('Scan the QR code with your WhatsApp app');

            // Wait for successful login
            await this.page.waitForSelector('[data-testid="chat-list"]', { timeout: 0 });
            console.log('Successfully logged in to WhatsApp Web');

            // Set up message observer
            await this.page.evaluate(() => {
                const observer = new MutationObserver((mutations) => {
                    mutations.forEach((mutation) => {
                        if (mutation.type === 'childList') {
                            const messages = document.querySelectorAll('[data-testid="msg-container"]');
                            messages.forEach((msg) => {
                                const text = msg.textContent;
                                const time = msg.querySelector('time')?.textContent;
                                const from = msg.closest('[data-testid="conversation-panel"]')?.querySelector('[data-testid="conversation-info-header"]')?.textContent;
                                
                                if (text && time && from) {
                                    window.postMessage({
                                        type: 'newMessage',
                                        data: { from, text, time }
                                    }, '*');
                                }
                            });
                        }
                    });
                });

                observer.observe(document.body, {
                    childList: true,
                    subtree: true
                });
            });

            // Listen for new messages
            this.page.on('console', async (msg) => {
                if (msg.type() === 'log' && msg.text().startsWith('newMessage:')) {
                    const data = JSON.parse(msg.text().slice(11));
                    this.broadcast({
                        type: 'message',
                        data: {
                            from: data.from,
                            content: data.text,
                            time: Date.now()
                        }
                    });
                }
            });
        }
    }

    async sendMessage(to, content) {
        try {
            // Search for contact
            await this.page.click('[data-testid="chat-list-search"]');
            await this.page.type('[data-testid="chat-list-search"]', to);
            await this.page.waitForTimeout(1000);

            // Click on the first chat that matches
            const chatSelector = `[title*="${to}"]`;
            await this.page.waitForSelector(chatSelector);
            await this.page.click(chatSelector);

            // Type and send message
            await this.page.waitForSelector('[data-testid="conversation-compose-box-input"]');
            await this.page.type('[data-testid="conversation-compose-box-input"]', content);
            await this.page.keyboard.press('Enter');

            return true;
        } catch (error) {
            console.error('Error sending message:', error);
            return false;
        }
    }

    broadcast(data) {
        const message = JSON.stringify(data);
        for (const client of this.clients) {
            if (client.readyState === WebSocket.OPEN) {
                client.send(message);
            }
        }
    }

    async stop() {
        if (this.browser) {
            await this.browser.close();
        }
        if (this.wss) {
            this.wss.close();
        }
    }
}

// Start WhatsApp Web automation
const whatsapp = new WhatsAppWeb();
whatsapp.start().catch(console.error);

// Handle graceful shutdown
process.on('SIGINT', async () => {
    console.log('Shutting down...');
    await whatsapp.stop();
    process.exit(0);
}); 