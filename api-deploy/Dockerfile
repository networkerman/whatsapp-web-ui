# Use Node.js LTS version
FROM node:18-alpine

# Set working directory
WORKDIR /app

# Copy package files
COPY api-deploy/package*.json ./

# Install dependencies
RUN npm ci

# Copy source code and config
COPY api-deploy/src ./src
COPY api-deploy/tsconfig.json ./

# Build TypeScript
RUN npm run build

# Create non-root user
RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs

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
CMD ["node", "dist/index.js"] 