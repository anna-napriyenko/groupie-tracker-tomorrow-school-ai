.PHONY: help build run stop restart clean logs test docker-build docker-up docker-down docker-logs docker-restart

# Variables
APP_NAME=groupie-tracker
DOCKER_IMAGE=$(APP_NAME):latest
BINARY=bin/server

## help: Show this help message
help:
	@echo 'Usage:'
	@echo '  make <target>'
	@echo ''
	@echo 'Targets:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## build: Build the Go application
build:
	@echo "Building application..."
	@mkdir -p bin
	@go build -o $(BINARY) cmd/api/main.go
	@echo "Build complete: $(BINARY)"

## run: Run the application locally
run: build
	@echo "Starting application..."
	@./$(BINARY)

## test: Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

## clean: Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@go clean
	@echo "Clean complete"

## docker-build: Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) .
	@echo "Docker image built: $(DOCKER_IMAGE)"

## docker-up: Start application with Docker Compose
docker-up:
	@echo "Starting application with Docker Compose..."
	@docker-compose up -d
	@echo "Application started. Access at http://localhost:8000"

## docker-down: Stop and remove Docker containers
docker-down:
	@echo "Stopping Docker containers..."
	@docker-compose down
	@echo "Containers stopped"

## docker-restart: Restart Docker containers
docker-restart:
	@echo "Restarting Docker containers..."
	@docker-compose restart
	@echo "Containers restarted"

## docker-logs: Show Docker container logs
docker-logs:
	@docker-compose logs -f

## docker-rebuild: Rebuild and restart Docker containers
docker-rebuild:
	@echo "Rebuilding and restarting..."
	@docker-compose down
	@docker-compose up -d --build
	@echo "Rebuild complete"

## docker-clean: Remove Docker images and containers
docker-clean:
	@echo "Cleaning Docker resources..."
	@docker-compose down -v --rmi all
	@echo "Docker resources cleaned"

## health: Check application health
health:
	@curl -f http://localhost:8000/health || echo "Health check failed"

## dev: Run application in development mode
dev:
	@echo "Starting development server..."
	@go run cmd/api/main.go
