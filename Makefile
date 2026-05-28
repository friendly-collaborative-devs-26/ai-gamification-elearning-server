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
tidy: depcheck
	@echo "$(GREEN)▶ tidying go.mod and go.sum...$(RESET)"
	@go mod tidy
	@echo "$(GREEN)✔ go.mod tidy complete$(RESET)"

.PHONY: depcheck
depcheck:
	@echo "$(GREEN)▶ checking for dependency issues...$(RESET)"
	@go mod verify
	@echo "$(GREEN)▶ checking minimum release age...$(RESET)"
	@go run ./cmd/depcheck
	@echo "$(GREEN)✔ dependency check complete$(RESET)"

.PHONY: depcheck-warn
depcheck-warn:
	@echo "$(YELLOW)▶ checking minimum release age (warn only)...$(RESET)"
	@go mod verify
	@DEPCHECK_MODE=warn go run ./cmd/depcheck
	@echo "$(YELLOW)✔ age check complete (warnings only)$(RESET)"

.PHONY: ci
ci: depcheck build
