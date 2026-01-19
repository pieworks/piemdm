# Getting Started

To set up the development environment, you can use **Docker Compose** (recommended for quick start) or set up the environment **manually**.

## Prerequisites

Ensure you have the following installed:
- **Go**: 1.24.12+
- **Node.js**: 20+
- **MySQL**: 8.0+
- **Redis**: 6+

## Quick Start (Docker Compose)

Start all services in the background:

```bash
cp deploy/.env.example .env
# Edit .env to set MySQL DSN and Redis address
vim .env

docker-compose -f deploy/docker-compose.yml up -d
```

This will build and start the **backend API** and **frontend** containers.

> [!IMPORTANT]
> This Docker Compose configuration **does not** include MySQL and Redis. You must have them running externally and configured via environment variables.

- **Frontend**: `http://localhost:8081`
- **Backend API**: `http://localhost:8787`
- **API Documentation**: `http://localhost:8787/swagger/index.html`

## Manual Setup

### 1. Backend

The backend is a Go application located in the `backend/` directory.

```bash
cd backend

# Install dependencies
go mod tidy

# Copy configuration file
cp config/local.yml.example config/local.yml
# Edit config/local.yml to configure database and other settings

# Run database migrations
go run cmd/migration/main.go

# Start backend service
go run cmd/server/main.go
```

### 2. Frontend

The frontend is a Vue.js application located in the `frontend/` directory.

```bash
cd frontend

# Install dependencies
npm install
# or
pnpm install

# Start development server
npm run dev
# or
pnpm dev
```

### 3. Database Setup

If running manually, ensure your database is created:

```sql
CREATE DATABASE piemdm CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

## Configuration

### Environment Variables (.env)

Create a `.env` file in the root directory for global settings:

```env
# Database
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=piemdm

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# Application
APP_PORT=8787
```
