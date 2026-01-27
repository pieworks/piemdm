# Getting Started

This guide will help you set up the **PieMDM** development environment.

## Prerequisites

Ensure you have the following installed on your local machine:

- [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)
- [Go](https://go.dev/) (>= 1.20) for backend development
- [Node.js](https://nodejs.org/) (>= 20.10.0) & [pnpm](https://pnpm.io/) (>= 9.0.0) for frontend development

## Quick Start (Docker)

The easiest way to launch the application stack is using Docker Compose. This will start the API and Web services. (Requires external MySQL & Redis)
最简单的启动应用栈的方法是使用 Docker Compose。这将启动 API 和 Web 服务。（需外部 MySQL 和 Redis 支持）

```bash
# Enter the deployment directory
cd deploy
# Start all services in the background
docker-compose up -d
```

Once started, you can access:
- **Web Interface**: [http://localhost:80](http://localhost:80)
- **API Server**: [http://localhost:8787](http://localhost:8787)

## Manual Setup

If you prefer to run services individually for development:

### 1. Database & Cache

You still need MySQL and Redis. You can use Docker to spin them up:

```bash
docker-compose up -d mysql redis
```

### 2. Backend (API)
The backend is developed in Go, based on the Gin framework.
Navigate to the `api` directory:

```bash
cd api
```

Install dependencies and run the server (with hot-reload enabled via `air`):

```bash
make dev
```

Common commands:
- `make test`: Run unit tests
- `make lint`: Run code linting
- `make wire`: Generate dependency injection code

### 3. Frontend (Web)
The frontend is developed in Vue 3.
Navigate to the `web` directory:

```bash
cd web
```

Install dependencies and start the development server:

```bash
pnpm install
pnpm dev
```

If you need to access it from other devices, pass the `--host` parameter:

```bash
pnpm dev --host
```

## Project Structure

- **`api/`**: Backend application (Go, Gin, GORM)
- **`web/`**: Frontend application (Vue 3, Vite, Bootstrap 5)
- **`docs/`**: Project documentation
