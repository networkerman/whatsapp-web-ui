#!/bin/bash

# Create directory for MCP server
mkdir -p whatsapp-mcp-server

# Copy files from original repository
cp -r ../../../whatsapp-mcp-server/* whatsapp-mcp-server/

# Install Python dependencies
cd whatsapp-mcp-server
pip install -r requirements.txt
cd ..

echo "MCP server setup complete!" 