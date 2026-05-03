APP_NAME    := ai-gamification-elearning-server
BINARY_DIR  := tmp
BINARY      := $(BINARY_DIR)/main
CMD_PATH    := ./cmd/main.go

GREEN  := \033[1;32m
YELLOW := \033[1;33m
RESET  := \033[0m

.PHONY: run
run:
	@echo "$(GREEN)▶ starting $(APP_NAME) in development mode...$(RESET)"
	@APP_ENV=development go run $(CMD_PATH)

.PHONY: run/watch
run/watch:
	@command -v air >/dev/null 2>&1 || { echo "Install air: go install github.com/air-verse/air@latest"; exit 1; }
	@echo "$(GREEN)▶ starting with live-reload (air)...$(RESET)"
	@air

.PHONY: build
build:
	@echo "$(GREEN)▶ building $(BINARY)...$(RESET)"
	@mkdir -p $(BINARY_DIR)
	@go build -o $(BINARY) $(CMD_PATH)
	@echo "$(GREEN)✔ build complete: $(BINARY)$(RESET)"

.PHONY: clean
clean:
	@echo "$(YELLOW)▶ cleaning build artifacts...$(RESET)"
	@rm -rf $(BINARY_DIR)
	@echo "$(GREEN)✔ clean complete$(RESET)"

.PHONY: tidy
tidy:
	@echo "$(GREEN)▶ tidying go.mod and go.sum...$(RESET)"
	@go mod tidy
	@echo "$(GREEN)✔ go.mod tidy complete$(RESET)"
