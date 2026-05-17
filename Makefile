# Makefile for go-clean-api
.DEFAULT_GOAL := help

.PHONY: help
help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

# ~~~ Code Actions ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.PHONY: lint
lint: ## Runs golangci-lint with predefined configuration
	@echo "> Applying linter"
	@go tool golangci-lint version
	@go tool golangci-lint run -c .golangci.yml ./...

.PHONY: format
format: ## Runs code formatter
	@echo "> Applying formatter"
	go tool golangci-lint fmt -c .golangci.yml ./...

.PHONY: deps
deps: ## Runs deps tidy + verify + check dependencies
	@echo "> Applying formatter"
	go mod tidy && go mod verify
	go tool govulncheck ./...

.PHONY: generate
generate: ## Generate go generte files
	@echo "> Generate go files"
	go generate ./...

SCHEMA_DIR := ./migrations
BOILER_SCHEMA := $(SCHEMA_DIR)/schema.sql.boiler
AUTOGEN_GO_DIR := ./internal/infra/mariadb/models
COMPOSE_FILE := deployments/docker-compose.yml
COMPOSE := docker compose -f $(COMPOSE_FILE)
COMPOSE_BOILER := $(COMPOSE) --profile boiler

.PHONY: generate.boilerplate
generate.boilerplate: ## Generate boilerplate db models (from schema.sql.boiler)
	@echo "> Generate boilerplate db models"
	@test -f $(BOILER_SCHEMA) || (echo "missing $(BOILER_SCHEMA)"; exit 1)
	mkdir -p $(AUTOGEN_GO_DIR)
	rm -f $(AUTOGEN_GO_DIR)/*.go
	$(COMPOSE_BOILER) up -d --wait mariadb-boiler
	PATH="$$(dirname $$(go tool -n sqlboiler-mysql)):$$PATH" \
		go tool sqlboiler mysql --config sqlboiler.toml --output $(AUTOGEN_GO_DIR) --no-tests --wipe
	$(COMPOSE_BOILER) down

.PHONY: generate.boilerplate.reset
generate.boilerplate.reset: ## Recreate generate DB volume then run generate.boilerplate
	$(COMPOSE_BOILER) down -v
	$(MAKE) generate.boilerplate

TESTS_ARGS := --format testname --jsonfile ./tmp/gotestsum.json.out
TESTS_ARGS += --max-fails 2
TESTS_ARGS += -- ./...
TESTS_ARGS += -test.parallel 2
TESTS_ARGS += -test.count    1
TESTS_ARGS += -test.failfast
TESTS_ARGS += -test.coverprofile ./tmp/coverage.out
TESTS_ARGS += -test.timeout 5s
TESTS_ARGS += -race

.PHONY: test
test: ## Runs go unittest (using gotestsum)
	@go tool gotestsum $(TESTS_ARGS)

.PHONY: clean
clean: ## Clean go module cache
	@go clean -cache
	@go clean -modcache

# ~~~ Development Environment ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.PHONY: db-up
db-up: ## Start dev MariaDB (deployments/docker-compose.yml)
	$(COMPOSE) up -d mariadb

.PHONY: compose-up
compose-up: ## Start API + MariaDB with build
	$(COMPOSE) up --build

.PHONY: compose-down
compose-down: ## Stop API + MariaDB with build
	$(COMPOSE) down --remove-orphans

.PHONY: run
run: ## Starts AIR ( Continuous Development restapi server).
	@go tool air
