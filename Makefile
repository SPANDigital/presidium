FILENAME=presidium
DOCSDIR=docs
.DEFAULT_GOAL=help
.PHONY: build test dist clean fmt vet tidy coverage_report help

help: ## Display available targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-18s %s\n", $$1, $$2}'

prepare-themes: ## Prepare themes for embedding (rename go.mod to make embeddable)
	@echo "Preparing themes for embedding..."
	@for theme in themes/presidium-*; do \
		if [ -f "$$theme/go.mod" ]; then \
			echo "  Renaming $$theme/go.mod -> go.mod.tmpl"; \
			mv "$$theme/go.mod" "$$theme/go.mod.tmpl"; \
		fi; \
		if [ -f "$$theme/go.sum" ]; then \
			echo "  Renaming $$theme/go.sum -> go.sum.tmpl"; \
			mv "$$theme/go.sum" "$$theme/go.sum.tmpl"; \
		fi; \
	done

restore-themes: ## Restore theme go.mod files to original names
	@echo "Restoring theme go.mod files..."
	@for theme in themes/presidium-*; do \
		if [ -f "$$theme/go.mod.tmpl" ]; then \
			echo "  Renaming $$theme/go.mod.tmpl -> go.mod"; \
			mv "$$theme/go.mod.tmpl" "$$theme/go.mod"; \
		fi; \
		if [ -f "$$theme/go.sum.tmpl" ]; then \
			echo "  Renaming $$theme/go.sum.tmpl -> go.sum"; \
			mv "$$theme/go.sum.tmpl" "$$theme/go.sum"; \
		fi; \
	done

build: prepare-themes ## Build the presidium binary
	go build -tags extended -o $(FILENAME) .

test: prepare-themes ## Run tests with coverage
	@mkdir -p reports
	go test -race -timeout 120s ./... -coverprofile=reports/tests-cov.out

fmt: ## Format Go source files
	go fmt ./...

vet: prepare-themes ## Run go vet
	go vet ./...

tidy: ## Tidy and verify module dependencies
	go mod tidy && go mod verify

clean: restore-themes ## Remove build artifacts and restore themes
	rm -fr "dist" "$(FILENAME)" "presidium-test"

coverage_report: ## Open coverage report in browser
	@go tool cover -html=reports/tests-cov.out

dist: prepare-themes ## Build distribution binary
	mkdir -p "dist"
	go build -trimpath -o "dist/presidium" --tags extended

checks: clean tidy fmt vet lint test build

serve-docs:
	cd $(DOCSDIR) && make serve

lint: prepare-themes ## Run golangci-lint
	golangci-lint run --timeout 10m