BINARY_NAME=pb

GO=go

OUTPUT_DIR=bin

build:
	@mkdir -p $(OUTPUT_DIR)
	@$(GO) build -o $(OUTPUT_DIR)/$(BINARY_NAME) .

run:
	@$(GO) run . --config=config.yml

dev:
	@$(GO) run . --config=config.dev.yml --port=8082

docker:
	@docker compose down -v
	@docker rmi pulsebridge -f || true
	@docker compose up -d
	@echo "Docker image built and containers started."

test:
	@$(GO) test ./...