SCRIPT_DIR = "scripts"

.PHONY: format-and-lint run-unit-tests run-integration-tests \
        spin-up-integration-test-docker-containers \
        spin-up-docker-containers shut-down-docker-containers help

# Help target to list all available targets
help:
	@echo "Available Makefile targets:"
	@echo "  format-and-lint                     - Run the format and linting script"
	@echo "  run-unit-tests                      - Run the unit tests"
	@echo "  run-integration-tests               - Run the integration tests"
	@echo "  spin-up-integration-test-docker-containers - Spin up Docker containers for integration tests (Postgres, Azure Blob Storage)"
	@echo "  spin-up-docker-containers           - Spin up Docker containers with internal containerized applications"
	@echo "  shut-down-docker-containers         - Shut down the application Docker containers"

# Run the format and lint script
format-and-lint:
	@cd $(SCRIPT_DIR) && ./format-and-lint.sh

# Run unit tests
run-unit-tests:
	@cd $(SCRIPT_DIR) && ./run-test.sh -u

# Run integration tests
run-integration-tests:
	@cd $(SCRIPT_DIR) && ./run-test.sh -i

# Spin up Docker containers for integration tests
spin-up-integration-test-docker-containers:
	docker-compose up -d postgres azure-blob-storage

# Spin up Docker containers with internal containerized applications
spin-up-docker-containers:
	docker-compose up -d --build

# Shut down Docker containers with internal containerized applications
shut-down-docker-containers:
	docker-compose down -v
