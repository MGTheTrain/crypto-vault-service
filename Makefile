SCRIPT_DIR = "scripts"

COVERAGE_OUT_FILE=coverage.out
COVERAGE_HTML_FILE=coverage.html
MIN_COVERAGE=80.0

# Help target to list all available targets
help:
	@echo "Available Makefile targets:"
	@echo "  format-and-lint                     		- Run the format and linting script"
	@echo "  lint-results			                    - Write golang-ci lint findings to a linter-findings.txt file"
	@echo "  run-unit-tests                      		- Run the unit tests"
	@echo "  run-integration-tests               		- Run the integration tests"
	@echo "  run-unit-and-integration-tests             - Run the unit and integration tests"
	@echo "  check-coverage                             - Run the unit and integration tests and check if code coverage of min 80 percent is achieved"
	@echo "  run-api-tests             					- Run the api tests"
	@echo "  spin-up-integration-test-docker-containers - Spin up Docker containers for integration tests (Postgres, Azure Blob Storage)"
	@echo "  spin-up-docker-containers           		- Spin up Docker containers with internal containerized applications"
	@echo "  shut-down-docker-containers         		- Shut down the application Docker containers"
	@echo "  generate-swagger-docs         				- Convert Go annotations to Swagger Documentation 2.0"
	@echo "  remove-artifacts         			 	    - Remove artifacts"

format-and-lint:
	@cd $(SCRIPT_DIR) && ./format-and-lint.sh

lint-results:
	@echo "Running golangci-lint..."
	@golangci-lint run | sed 's/^/- /' > linter-findings.txt
	@echo "Linting results written to linter-findings.txt"

run-unit-tests:
	@echo "Running unit tests..."
	@go test ./internal/... --tags="integration" -cover

run-integration-tests:
	@echo "Running integration tests..."
	@go test ./internal/... --tags="integration" -cover

run-unit-and-integration-tests:
	@echo "Running unit and integration tests... Generating $(COVERAGE_HTML_FILE) file..."
	@go test ./internal/... --tags="unit integration" -cover -coverprofile=$(COVERAGE_OUT_FILE) -v
	@go tool cover -html=$(COVERAGE_OUT_FILE) -o $(COVERAGE_HTML_FILE)

check-coverage: run-unit-and-integration-tests
	@echo "Checking if coverage meets minimum threshold ($(MIN_COVERAGE)%)..."
	@total_coverage=$$(go tool cover -func=$(COVERAGE_OUT_FILE) | grep total | awk '{print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$total_coverage < $(MIN_COVERAGE)" | bc) -eq 1 ]; then \
		echo "❌ Code coverage ($$total_coverage%) is below the required $(MIN_COVERAGE)% threshold"; \
		exit 1; \
	else \
		echo "✅ Code coverage check passed: $$total_coverage%"; \
	fi

run-api-tests:
	@cd $(SCRIPT_DIR) && echo "TODO(MGTheTrain): Invoke API tests"

run-e2e-tests:
	@echo "Running e2e tests..."
	@go test ./test/... --tags="e2e" -cover

spin-up-integration-test-docker-containers:
	docker compose up -d postgres azure-blob-storage

spin-up-docker-containers:
	docker compose up -d --build

shut-down-docker-containers:
	docker compose down -v

generate-swagger-docs:
	@echo "Generating Swagger docs..."
	@swag init -g cmd/crypto-vault-rest-service/crypto_vault_service.go -o cmd/crypto-vault-rest-service/docs

generate-grpc-files:
	@cd $(SCRIPT_DIR) && ./generate-grpc-files.sh

remove-artifacts:
	@echo "Removing artifacts..."
	@rm coverage.* linter-findings.*
