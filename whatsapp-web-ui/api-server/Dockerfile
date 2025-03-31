# Base stage for common dependencies
FROM node:18-alpine AS base
RUN apk add --no-cache python3 py3-pip

# Python stage for MCP server dependencies
FROM base AS python
RUN pip3 install --no-cache-dir requests

# Dependencies stage
FROM base AS deps
WORKDIR /app
COPY package*.json ./
RUN npm ci

# Builder stage
FROM deps AS builder
WORKDIR /app
COPY . .
RUN npm run build

# Runner stage
FROM base AS runner
WORKDIR /app

# Create non-root user
RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs

# Copy built files from builder
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/package*.json ./
COPY --from=builder /app/node_modules ./node_modules

# Copy MCP server files from the parent directory
COPY ../whatsapp-mcp-server ./whatsapp-mcp-server

# Install Python dependencies for MCP server
COPY --from=python /usr/local/lib/python3.9/site-packages /usr/local/lib/python3.9/site-packages

# Set environment variables
ENV NODE_ENV=production
ENV PORT=3001

# Set user permissions
RUN chown -R nextjs:nodejs /app

# Switch to non-root user
USER nextjs

# Expose port
EXPOSE 3001

# Start the server
CMD ["npm", "start"] 