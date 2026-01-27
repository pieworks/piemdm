# What is PieMDM?

PieMDM is a full-stack Master Data Management (MDM) platform designed to help enterprises centrally manage and maintain critical business data.

## Project Overview

PieMDM adopts a modern separation of frontend and backend architecture, providing an efficient, stable, and easily scalable master data management solution.

- **Backend**: Built with Go, utilizing the Gin framework to provide high-performance API services.
- **Frontend**: Built with Vue.js 3 and Vite, providing a responsive and user-friendly web interface.
- **Containerization**: Supports Docker and Docker Compose for easy development and deployment.

## Tech Stack

### Backend (API)

The backend code is located in the `api/` directory. Key technology choices include:

- **Language**: Go
- **Web Framework**: Gin
- **Database ORM**: GORM (MySQL)
- **Cache**: Redis
- **Dependency Injection**: google/wire
- **Testing**: testify, go-sqlmock

### Frontend (Web)

The frontend code is located in the `web/` directory. Key technology choices include:

- **Framework**: Vue.js 3
- **Build Tool**: Vite
- **State Management**: Pinia
- **Routing**: Vue Router
- **UI Framework**: Bootstrap 5
- **Package Manager**: pnpm

## Quick Start

PieMDM recommends using Docker Compose for a quick start:

```bash
# Enter the deployment directory
cd deploy
# Start all services in the background
docker-compose up -d
```

Once started, you can access:
- **API Service**: `http://localhost:8787`
- **Frontend Page**: `http://localhost:80`

## Learn More

- [Quick Start Guide](./getting-started)
- [API Reference](../reference/open-api)