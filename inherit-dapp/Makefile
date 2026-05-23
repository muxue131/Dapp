.PHONY: all build build-chain build-backend build-frontend dev dev-chain dev-backend dev-frontend test clean docker-up docker-down

# === Build ===

all: build

build: build-chain build-backend build-frontend

build-chain:
	@echo "Building chain..."
	cd chain && go build -o bin/legacyd ./cmd/legacyd/

build-backend:
	@echo "Building backend..."
	cd backend && go build -o bin/legacy-api ./main.go

build-frontend:
	@echo "Building frontend..."
	cd frontend && npm install && npm run build

# === Development ===

dev: dev-chain dev-backend dev-frontend

dev-chain:
	@echo "Starting chain node..."
	cd chain && go run ./cmd/legacyd/ serve

dev-backend:
	@echo "Starting backend server..."
	cd backend && go run ./main.go

dev-frontend:
	@echo "Starting frontend dev server..."
	cd frontend && npm run dev

# === Testing ===

test: test-chain test-backend

test-chain:
	@echo "Testing chain..."
	cd chain && go test ./... -v

test-backend:
	@echo "Testing backend..."
	cd backend && go test ./... -v

# === Docker ===

docker-up:
	docker-compose up -d --build

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

# === Clean ===

clean:
	rm -rf chain/bin/ backend/bin/ frontend/dist/ frontend/node_modules/

# === Lint ===

lint: lint-go lint-ts

lint-go:
	@echo "Linting Go code..."
	cd chain && golangci-lint run ./...
	cd backend && golangci-lint run ./...

lint-ts:
	@echo "Linting TypeScript code..."
	cd frontend && npm run lint

# === Genesis ===

init-chain:
	@echo "Initializing chain..."
	cd chain && go run ./cmd/legacyd/ init legacy-node

# === Install dependencies ===

deps:
	cd backend && go mod tidy
	cd frontend && npm install

# === Help ===

help:
	@echo "Legacy DApp - Available commands:"
	@echo ""
	@echo "  make build          - Build all components"
	@echo "  make build-chain    - Build chain binary"
	@echo "  make build-backend  - Build backend binary"
	@echo "  make build-frontend - Build frontend"
	@echo ""
	@echo "  make dev            - Start all dev servers"
	@echo "  make dev-chain      - Start chain node"
	@echo "  make dev-backend    - Start backend server"
	@echo "  make dev-frontend   - Start frontend dev server"
	@echo ""
	@echo "  make test           - Run all tests"
	@echo "  make test-chain     - Run chain tests"
	@echo "  make test-backend   - Run backend tests"
	@echo ""
	@echo "  make docker-up      - Start with Docker Compose"
	@echo "  make docker-down    - Stop Docker Compose"
	@echo ""
	@echo "  make clean          - Clean build artifacts"
	@echo "  make deps           - Install dependencies"
	@echo "  make init-chain     - Initialize chain node"
