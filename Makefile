SCRIPT_DIR = "scripts"

.PHONY: format-and-lint run-unit-tests run-integration-tests

format-and-lint:
	@cd $(SCRIPT_DIR) && ./format-and-lint.sh

run-unit-tests:
	@cd $(SCRIPT_DIR) && ./run-test.sh -u

run-integration-tests:
	@cd $(SCRIPT_DIR) && ./run-test.sh -i
