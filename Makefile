SCRIPT_DIR = "scripts"

# Help target to list all available targets
help:
	@echo "Available Makefile targets:"
	@echo "  format-and-lint                     		- Run the format and linting script"
	@echo "  run-unit-tests                      		- Run the unit tests"
	@echo "  run-integration-tests               		- Run the integration tests"
	@echo "  run-unit-and-integration-tests             - Run the unit and integration tests"
	@echo "  spin-up-integration-test-docker-containers - Spin up Docker containers for integration tests (Postgres, Azure Blob Storage)"
	@echo "  spin-up-docker-containers           		- Spin up Docker containers with internal containerized applications"
	@echo "  shut-down-docker-containers         		- Shut down the application Docker containers"
	@echo "  generate-swagger-docs         				- Convert Go annotations to Swagger Documentation 2.0"

format-and-lint:
	@cd $(SCRIPT_DIR) && ./format-and-lint.sh

run-unit-tests:
	@cd $(SCRIPT_DIR) && ./run-test.sh -u

run-integration-tests:
	@cd $(SCRIPT_DIR) && ./run-test.sh -i

run-unit-and-integration-tests:
	@cd $(SCRIPT_DIR) && ./run-test.sh -a

spin-up-integration-test-docker-containers:
	docker-compose up -d postgres azure-blob-storage

spin-up-docker-containers:
	docker-compose up -d --build

shut-down-docker-containers:
	docker-compose down -v

generate-swagger-docs:
	@cd $(SCRIPT_DIR) && ./generate-docs.sh

generate-grpc-files:
	@cd $(SCRIPT_DIR) && ./generate-grpc-files.sh
