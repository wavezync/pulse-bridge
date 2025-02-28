BINARY_NAME=pb

GO=go

OUTPUT_DIR=bin

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(OUTPUT_DIR)
	@$(GO) build -o $(OUTPUT_DIR)/$(BINARY_NAME) ./cmd/pulse-bridge/main.go

run:
	@echo "Running $(BINARY_NAME)..."
	@cd $(CURDIR) && $(GO) run ./cmd/pulse-bridge/main.go