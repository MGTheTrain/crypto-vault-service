SCRIPT_DIR = "scripts"

.PHONY: lint run-unit-tests run-integration-tests

lint:
	@cd $(SCRIPT_DIR) && ./format-and-lint.sh

run-unit-tests:
	@cd $(SCRIPT_DIR) && ./run-test.sh -u

run-integration-tests:
	@cd $(SCRIPT_DIR) && ./run-test.sh -i
