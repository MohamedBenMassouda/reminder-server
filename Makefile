# Set the path to the migration file
MIGRATE_CMD=go run ./cmd/migrate/main.go

.PHONY: run
run:  ## Run the application
	go run ./cmd/server/main.go

.PHONY: help
help:  ## Show help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: create
create: ## Create a new migration file
	@if [ -z "$(name)" ]; then \
		echo "Please specify a migration name with 'make create name=<migration_name>'"; \
	else \
		$(MIGRATE_CMD) create $(name); \
	fi

# Apply all migrations
.PHONY: up
up:  ## Apply all migrations
	$(MIGRATE_CMD) up

# Roll back one migration
.PHONY: down
down:  ## Roll back one migration
	$(MIGRATE_CMD) down

# Check migrations status
.PHONY: status
status:  ## Check migrations status
	$(MIGRATE_CMD) status

# Reset database (down all migrations then up all migrations)
.PHONY: reset
reset:  ## Reset database (down all migrations then up all migrations)
	$(MIGRATE_CMD) reset

down-to: ## Roll back to a specific migration
	@if [ -z "$(version)" ]; then \
		echo "Please specify a migration version with 'make down-to version=<migration_version>'"; \
	else \
		$(MIGRATE_CMD) down-to $(version); \
	fi

.PHONY: build
build:  ## Build the application
	go build -o bin/server ./cmd/server/main.go

.PHONY: test
test:  ## Run tests
	go test -v ./...

