# .RECIPEPREFIX := $(.RECIPEPREFIX)<space>
TESTCOVERAGE_THRESHOLD=0

default: help

.PHONY: help
help: ## Shows help for each of the Makefile recipes.
	@echo 'Commands:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-32s\033[0m %s\n", $$1, $$2}'

.PHONY: dep
dep: ## Downloads dependencies.
	go mod download
	go mod tidy

.PHONY: init-tools
init-tools: ## Initializes tools.
	command -v pre-commit >/dev/null || brew install pre-commit
	[ -f .git/hooks/pre-commit ] || pre-commit install
	command -v act >/dev/null || brew install act
	go tool -n air >/dev/null || go get -tool github.com/air-verse/air@latest

.PHONY: init-generators
init-generators: ## Initializes generators.
	go tool -n sqlc >/dev/null || go get -tool github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go tool -n mockery > /dev/null || go get -tool github.com/vektra/mockery/v2@latest
	go tool -n stringer >/dev/null || go get -tool golang.org/x/tools/cmd/stringer@latest
	go tool -n gcov2lcov >/dev/null || go get -tool github.com/jandelgado/gcov2lcov@latest

.PHONY: init-checkers
init-checkers: ## Initializes checkers.
	go tool -n golangci-lint >/dev/null || go get -tool github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	go tool -n betteralign >/dev/null || go get -tool github.com/dkorunic/betteralign/cmd/betteralign@latest
	go tool -n govulncheck >/dev/null || go get -tool golang.org/x/vuln/cmd/govulncheck@latest

.PHONY: init
init: init-tools init-generators init-checkers dep # Initializes the project.
	# cp -n .env.example .env || true

.PHONY: generate
generate: ## Runs auto-generated code generation tools.
	go generate ./...

.PHONY: migrate-up
migrate-up: ## Runs the migration up command.
	go run ./cmd/migrate/ default up

.PHONY: migrate-down
migrate-down: ## Runs the migration down command.
	go run ./cmd/migrate/ default down

.PHONY: build
build: ## Builds the entire codebase.
	go build -v ./...

.PHONY: clean
clean: ## Cleans the entire codebase.
	go clean

.PHONY: dev
dev: ## Runs the service in development mode.
	go tool air --build.bin "./tmp/serve" --build.cmd "go build -o ./tmp/serve ./cmd/serve/"

.PHONY: run
run: ## Runs the service.
	go run ./cmd/serve/

.PHONY: test
test: ## Runs the tests.
	go test -failfast -race -count 1 ./...

.PHONY: test-cov
test-cov: ## Runs the tests with coverage.
	go test -failfast -race -count 1 -coverpkg=./... -coverprofile=${TMPDIR}cov_profile.out ./...
	# go tool gcov2lcov -infile ${TMPDIR}cov_profile.out -outfile ./cov_profile.lcov

.PHONY: test-view-html
test-view-html: ## Views the test coverage in HTML.
	go tool cover -html ${TMPDIR}cov_profile.out -o ${TMPDIR}cov_profile.html
	open ${TMPDIR}cov_profile.html

.PHONY: test-ci
test-ci: test-cov # Runs the tests with coverage and check if it's above the threshold.
	$(eval ACTUAL_COVERAGE := $(shell go tool cover -func=${TMPDIR}cov_profile.out | grep total | grep -Eo '[0-9]+\.[0-9]+'))

	@echo "Quality Gate: checking test coverage is above threshold..."
	@echo "Threshold             : $(TESTCOVERAGE_THRESHOLD) %"
	@echo "Current test coverage : $(ACTUAL_COVERAGE) %"

	@if [ "$(shell echo "$(ACTUAL_COVERAGE) < $(TESTCOVERAGE_THRESHOLD)" | bc -l)" -eq 1 ]; then \
    echo "Current test coverage is below threshold. Please add more unit tests or adjust threshold to a lower value."; \
    echo "Failed"; \
    exit 1; \
  else \
    echo "OK"; \
  fi

.PHONY: lint
lint: ## Runs the linting command.
	go tool golangci-lint run ./...

.PHONY: check
check: ## Runs static analysis tools.
	go run golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest -test ./...
	go tool govulncheck ./...
	go tool betteralign ./...
	go vet ./...

.PHONY: fix
fix: ## Fixes code formatting and alignment.
	go run golang.org/x/tools/gopls/internal/analysis/modernize/cmd/modernize@latest -fix ./...
	go tool betteralign -apply ./...
	go fmt ./...

.PHONY: ok
ok: lint check test ## Checks if the code is ok.

%:
	@:
