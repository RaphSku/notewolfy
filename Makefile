.PHONY: default
default: help

.PHONY: start
## Run the CLI tool
start: 
	@go run main.go

.PHONY: unit-tests
## Run Unit Tests
unit-tests:
	go test -coverprofile=coverage_unit.out -cover ./... -tags unit_test
	go tool cover -html coverage_unit.out -o coverage_unit.html

.PHONY: integration-tests
## Run Integration Tests
integration-tests:
	go test -coverprofile=coverage_int.out -cover ./... -tags integration_test
	go tool cover -html coverage_int.out -o coverage_int.html

.PHONY: e2e-tests
## Run E2E Tests
e2e-tests:
	go test -coverprofile=coverage_e2e.out -cover ./... -tags e2e_test
	go tool cover -html coverage_e2e.out -o coverage_e2e.html

.PHONY: run-test-suite
## Run Complete Test Suite
run-test-suite: unit-tests integration-tests e2e-tests

help:
	@echo "----------------------------------"
	@echo "Welcome to make! Enjoy the flight."
	@echo "Makefile - make [\033[38;5;154mtarget\033[0m]"
	@echo "----------------------------------"
	@echo
	@echo "Targets:"
	@awk '/^[a-zA-z\-_0-9%:\\]+/ { \
		description = match(descriptionLine, /^## (.*)/); \
		if (description) { \
			target = $$1; \
			description = substr(descriptionLine, RSTART + 3, RLENGTH); \
			gsub("\\\\", "", target); \
			gsub(":+$$", "", target); \
			printf "    \033[38;5;154m%-25s\033[0m %s\n", target, description; \
		} \
	} \
	{ descriptionLine = $$0 }' $(MAKEFILE_LIST)
	@printf "\n"
