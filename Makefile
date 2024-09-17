# ℹ Freely based on: https://gist.github.com/thomaspoignant/5b72d579bd5f311904d973652180c705

# Constants
TARGET_FOLDER = target

# Docker images
DOCKER_IMAGE_GOLANG    = golang:1.21-alpine3.17
DOCKER_IMAGE_GOLANG_CI = golangci/golangci-lint:v1.59

# Some colors
COLOR_GREEN  = $(shell tput -Txterm setaf 2)
COLOR_YELLOW = $(shell tput -Txterm setaf 3)
COLOR_WHITE  = $(shell tput -Txterm setaf 7)
COLOR_CYAN   = $(shell tput -Txterm setaf 6)
COLOR_RED    = $(shell tput -Txterm setaf 1)
COLOR_RESET  = $(shell tput -Txterm sgr0)

.PHONY: all
all: help

## Lint:
.PHONY: lint
lint: lint-go

.PHONY: lint-go
lint-go: ## Lint go source code
	@echo "${COLOR_CYAN}🔍 Inspecting go source code${COLOR_RESET}"
	@docker run --rm \
  		-v `pwd`:/app:ro \
  		-w /app \
  		${DOCKER_IMAGE_GOLANG_CI} \
  		golangci-lint run -v

## Format:
.PHONY: format
format: format-go ## Run all available formatters

.PHONY: format-go
format-go: ## Format go files
	@echo "${COLOR_CYAN}📐 Formatting go source code${COLOR_RESET}"
	@docker run --rm \
  		-v `pwd`:/app:rw \
  		-w /app \
  		${DOCKER_IMAGE_GOLANG} \
  		sh -c \
		"go install mvdan.cc/gofumpt@v0.6.0; gofumpt -w -l ."

## Build:
.PHONY: build
build: build-go

.PHONY: build-go
build-go:
	@echo "${COLOR_CYAN} 🏗️ Building project${COLOR_RESET}"
	@go build ./...

## Test:
.PHONY: test
test: test-go ## Pass all the tests

.PHONY: test-go
test-go: build ## Pass the test for the go source code
	@echo "${COLOR_CYAN} 🧪 Passing go tests${COLOR_RESET}"
	@mkdir $(TARGET_FOLDER)
	@go test -v -coverprofile $(TARGET_FOLDER)/coverage.txt ./...

## Clean:
.PHONY: clean
clean: ## Remove all the files from the target folder
	@echo "${COLOR_CYAN} 🗑 Cleaning folder $(TARGET_FOLDER)${COLOR_RESET}"
	@rm -rf $(TARGET_FOLDER)/


## Mock:
.PHONY: mock
mock: ## Generate all the mocks (for tests)
	@echo "${COLOR_CYAN} 🧱 Generating all the mocks${COLOR_RESET}"
	@go install go.uber.org/mock/mockgen@v0.4.0
	@mockgen -source=auth/proxy.go -package testutil -destination testutil/auth_mocks.go
	@mockgen -source=dataverse/client.go -mock_names TxClient=MockDataverseTxClient  -package testutil -destination testutil/dataverse_mocks.go
	@mockgen -source=credential/parser.go -package testutil -destination testutil/credential_mocks.go
	@mockgen -package testutil -destination testutil/dataverse_client_mocks.go -mock_names QueryClient=MockDataverseQueryClient github.com/axone-protocol/axone-contract-schema/go/dataverse-schema/v5 QueryClient
	@mockgen -package testutil -destination testutil/cognitarium_client_mocks.go -mock_names QueryClient=MockCognitariumQueryClient github.com/axone-protocol/axone-contract-schema/go/cognitarium-schema/v5 QueryClient
	@mockgen -package testutil -destination testutil/law_stone_client_mocks.go -mock_names QueryClient=MockLawStoneQueryClient github.com/axone-protocol/axone-contract-schema/go/law-stone-schema/v5 QueryClient
	@mockgen -package testutil -destination testutil/auth_client_mocks.go -mock_names QueryClient=MockAuthQueryClient github.com/cosmos/cosmos-sdk/x/auth/types QueryClient
	@mockgen -package testutil -destination testutil/tx_service_mocks.go -mock_names ServiceClient=MockTxServiceClient github.com/cosmos/cosmos-sdk/types/tx ServiceClient
	@mockgen -source=credential/generate.go -package testutil -destination testutil/generate_mocks.go
	@mockgen -source=tx/transaction.go -package testutil -destination testutil/transaction_mocks.go
	@mockgen -source=tx/client.go -mock_names Client=MockTxClient -package testutil -destination testutil/tx_mocks.go
	@mockgen -source=keys/keyring.go -package testutil -destination testutil/keyring_mocks.go

## Help:
.PHONY: help
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${COLOR_YELLOW}make${COLOR_RESET} ${COLOR_GREEN}<target>${COLOR_RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${COLOR_YELLOW}%-20s${COLOR_GREEN}%s${COLOR_RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${COLOR_CYAN}%s${COLOR_RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)
	@echo ''
	@echo 'This Makefile depends on ${COLOR_CYAN}docker${COLOR_RESET}. To install it, please follow the instructions:'
	@echo '- for ${COLOR_YELLOW}macOS${COLOR_RESET}: https://docs.docker.com/docker-for-mac/install/'
	@echo '- for ${COLOR_YELLOW}Windows${COLOR_RESET}: https://docs.docker.com/docker-for-windows/install/'
	@echo '- for ${COLOR_YELLOW}Linux${COLOR_RESET}: https://docs.docker.com/engine/install/'
