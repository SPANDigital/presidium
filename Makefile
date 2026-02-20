FILENAME=presidium
DOCSDIR=docs
.DEFAULT_GOAL=help
.PHONY: build test dist clean fmt vet tidy coverage_report help

help: ## Display available targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-18s %s\n", $$1, $$2}'

build: ## Build the presidium binary
	go build -tags extended -o $(FILENAME) .

test: ## Run tests with coverage
	@mkdir -p reports
	go test -race -timeout 120s ./... -coverprofile=reports/tests-cov.out

fmt: ## Format Go source files
	go fmt ./...

vet: ## Run go vet
	go vet ./...

tidy: ## Tidy and verify module dependencies
	go mod tidy && go mod verify

clean: ## Remove build artifacts
	rm -fr "dist"

coverage_report: ## Open coverage report in browser
	@go tool cover -html=reports/tests-cov.out

dist: ## Build distribution binary
	mkdir -p "dist"
	go build -trimpath -o "dist/presidium" --tags extended

checks: clean tidy fmt vet lint test build

serve-docs:
	cd $(DOCSDIR) && make serve

lint:
	golangci-lint run --timeout 10m