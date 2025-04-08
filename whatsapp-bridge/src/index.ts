import WhatsAppClient from './whatsapp';

async function main() {
    const client = new WhatsAppClient(8080);
    
    try {
        await client.start();
        console.log('WhatsApp client started successfully');
    } catch (error) {
        console.error('Failed to start WhatsApp client:', error);
        process.exit(1);
    }

    // Handle graceful shutdown
    process.on('SIGINT', async () => {
        console.log('Shutting down...');
        await client.stop();
        process.exit(0);
    });
}

main(); 