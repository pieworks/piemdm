.PHONY: help dev dev-backend dev-frontend test-all clean build-all

# Default target: Show help
help:
	@echo "PieMDM Monorepo Management Commands"
	@echo ""
	@echo "Development Commands:"
	@echo "  make dev              - Start full development environment (Docker)"
	@echo "  make dev-backend      - Start backend development service only"
	@echo "  make dev-frontend     - Start frontend development service only"
	@echo ""
	@echo "Testing Commands:"
	@echo "  make test-all         - Run all tests"
	@echo "  make test-backend     - Run backend tests"
	@echo "  make test-frontend    - Run frontend tests"
	@echo ""
	@echo "Build Commands:"
	@echo "  make build-all        - Build all production images"
	@echo "  make clean            - Clean up environment"

# Start development environment
dev:
	docker-compose -f deploy/docker-compose.yml up -d
	@echo "✅ Development environment started"
	@echo "   Backend: http://localhost:8787"
	@echo "   Frontend: http://localhost:80"

# Start backend development only
dev-backend:
	cd backend && go run cmd/api/main.go

# Start frontend development only
dev-frontend:
	cd frontend && pnpm dev

# Run all tests
test-all: test-backend test-frontend

# Backend tests
test-backend:
	@echo "🧪 Running backend tests..."
	cd backend && go test -v ./...

# Frontend tests
test-frontend:
	@echo "🧪 Running frontend tests..."
	cd frontend && pnpm test

# Build all production images
build-all:
	@echo "🏗️  Building production images..."
	docker build -f backend/Dockerfile -t piemdm-backend:latest .
	docker build -f frontend/Dockerfile -t piemdm-frontend:latest .
	@echo "✅ Build complete"

# Clean up environment
clean:
	docker-compose -f deploy/docker-compose.yml down
	docker system prune -f
	@echo "✅ Environment cleaned"

# Start production environment
prod:
	docker-compose -f deploy/docker-compose.prod.yml up -d
	@echo "✅ Production environment started"

# Stop production environment
prod-down:
	docker-compose -f deploy/docker-compose.prod.yml down

# Generate API Client
gen-api: ## Generate Swagger docs and TypeScript client
	@echo "🔄 Generating API client..."
	./scripts/gen-api-client.sh
