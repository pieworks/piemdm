# Multi-stage build for Vue.js application
FROM node:20-alpine AS builder

# Set working directory
WORKDIR /app

# Copy workspace configs
COPY pnpm-lock.yaml pnpm-workspace.yaml ./

# Copy package definitions
COPY frontend/package.json ./frontend/
COPY shared/types/package.json ./shared/types/
COPY packages/workflow-vue/package.json ./packages/workflow-vue/

# Remove prepare script to prevent premature build during install
RUN node -e "const fs=require('fs'); const path='./shared/types/package.json'; const pkg=JSON.parse(fs.readFileSync(path)); delete pkg.scripts.prepare; fs.writeFileSync(path, JSON.stringify(pkg, null, 2));"

# Install dependencies (frozen-lockfile compatible)
RUN npm install -g pnpm && pnpm install --frozen-lockfile

# Copy source code
COPY frontend/ ./frontend/
COPY shared/types/ ./shared/types/
COPY packages/workflow-vue/ ./packages/workflow-vue/

# Build shared types first
WORKDIR /app/shared/types
RUN pnpm build

# Build workflow-vue
WORKDIR /app/packages/workflow-vue
RUN pnpm build

# Build frontend
WORKDIR /app/frontend
RUN pnpm build

# Production stage
FROM nginx:alpine

# Copy built assets from builder stage
COPY --from=builder /app/frontend/dist /usr/share/nginx/html

# Copy nginx configuration
COPY frontend/nginx.conf /etc/nginx/conf.d/default.conf

# Create nginx user and set permissions
RUN chown -R nginx:nginx /usr/share/nginx/html && \
    chmod -R 755 /usr/share/nginx/html

# Expose port
EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost/ || exit 1

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
