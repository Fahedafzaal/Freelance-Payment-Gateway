.PHONY: help build run test clean deploy test-api docker

# Default target
help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Development
build: ## Build the application
	go build -o bin/payment-gateway cmd/main.go

run: ## Run the application
	go run cmd/main.go

test: ## Run tests
	go test ./...

test-short: ## Run tests (short mode, skip integration tests)
	go test -short ./...

clean: ## Clean build artifacts
	rm -rf bin/

# Dependencies
deps: ## Download dependencies
	go mod download
	go mod tidy

# Smart Contract
compile-contract: ## Compile smart contract (run from parent directory)
	cd ../ && forge build

deploy-contract: ## Deploy smart contract to testnet
	./scripts/deploy.sh

# API Testing
test-api: ## Test API endpoints
	./scripts/test-api.sh

# Development setup
setup: ## Initial setup (copy env file, install deps)
	cp env.example .env
	@echo "📝 Please edit .env file with your configuration"
	go mod download

# Linting and formatting
fmt: ## Format Go code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

# Docker (optional)
docker-build: ## Build Docker image
	docker build -t freelance-payment-gateway .

docker-run: ## Run Docker container
	docker run -p 8081:8081 --env-file .env freelance-payment-gateway

# Production
build-prod: ## Build for production
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/payment-gateway cmd/main.go

# All-in-one commands
dev-setup: setup compile-contract ## Complete development setup
	@echo "✅ Development setup complete!"
	@echo "Next steps:"
	@echo "1. Edit .env file with your configuration"
	@echo "2. Run 'make deploy-contract' to deploy smart contract"
	@echo "3. Run 'make run' to start the server"

check: fmt vet test-short ## Run all checks (format, vet, test) 