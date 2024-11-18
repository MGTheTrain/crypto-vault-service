SCRIPT_DIR = "scripts"

.PHONY: format-and-lint run-unit-tests run-integration-tests \
        spin-up-integration-test-docker-containers \
        shut-down-integration-test-docker-containers \
        spin-up-docker-containers shut-down-docker-containers
		
format-and-lint:
	@cd $(SCRIPT_DIR) && ./format-and-lint.sh

run-unit-tests:
	@cd $(SCRIPT_DIR) && ./run-test.sh -u

run-integration-tests:
	@cd $(SCRIPT_DIR) && ./run-test.sh -i

spin-up-integration-test-docker-containers:
	docker-compose up -d postgres azure-blob-storage

shut-down-integration-test-docker-containers:
	docker-compose down postgres azure-blob-storage -v

spin-up-docker-containers:
	docker-compose up -d --build

shut-down-docker-containers:
	docker-compose down -v