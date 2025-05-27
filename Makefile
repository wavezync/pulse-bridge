BINARY_NAME=pb

GO=go

OUTPUT_DIR=bin

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(OUTPUT_DIR)
	@$(GO) build -o $(OUTPUT_DIR)/$(BINARY_NAME) .

run:
	@echo "Running $(BINARY_NAME)..."
	@$(GO) run . --config=config.yml

dev:
	@echo "Running in development mode..."
	@$(GO) run . --config=config.dev.yml --port=8082

docker:
	@echo "Building Docker image..."
	@docker compose down -v
	@docker rmi pulsebridge -f || true
	@docker build ./ -t pulsebridge --no-cache
	@docker compose up -d
	@echo "Docker image built and containers started."