BINARY_NAME=pulse-bridge

GO=go

OUTPUT_DIR=bin

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(OUTPUT_DIR)
	@$(GO) build -o $(OUTPUT_DIR)/$(BINARY_NAME) ./main.go

run:
	@echo "Running $(BINARY_NAME)..."
	@$(GO) run ./main.go