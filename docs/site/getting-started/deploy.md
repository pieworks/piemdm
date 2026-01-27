# Deployment Guide

PieMDM is designed to be easily deployed using containerization technologies. We recommend using Docker Compose for production deployment.

## Prerequisites

The deployment environment requires:

- [Docker](https://docs.docker.com/engine/install/)
- [Docker Compose](https://docs.docker.com/compose/install/)

## Deploy with Docker Compose

1. **Get Code**

   Clone the project code to the server:

   ```bash
   git clone https://github.com/pieteams/piemdm.git
   cd piemdm
   ```

2. **Configuration**

   Modify `docker-compose.yaml` or related configuration files (such as database passwords, port mappings, etc.) according to your environment needs.

3. **Start Services**

   Run in the project root directory:

   ```bash
   cd deploy
   docker-compose up -d
   ```

   This command will build and start the application services (API, Web).

4. **Verify**

   Use `docker-compose ps` to check service status.
   
   - Default Frontend Port: `80`
   - Default API Server Port: `8787`

## Manual Build & Deploy

If you need to manually build binaries and static resources:

### Backend

```bash
cd api
make build
# Run the generated binary
./bin/server
```

### Frontend

```bash
cd web
pnpm install
pnpm build
# Build artifacts are located in the web/dist directory, ready for hosting with Nginx or other web servers
```
