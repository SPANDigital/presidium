FILENAME=presidium
DOCSDIR=docs
.DEFAULT_GOAL=help
.PHONY: build test test-offline serve-offline dist clean fmt vet tidy coverage_report help update-themes prepare-themes lint checks serve-docs

help: ## Display available targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-18s %s\n", $$1, $$2}'

update-themes: ## Update theme submodules to their latest versions
	@echo "Updating theme submodules to latest..."
	@git submodule update --remote themes/presidium-styling-base themes/presidium-layouts-base themes/presidium-layouts-blog
	@echo "Theme submodules updated."

prepare-themes: ## Prepare themes for embedding (create zip bundle)
	@rm -f themes.zip
	@echo "Creating themes.zip bundle..."
	@cd themes && zip -r ../themes.zip presidium-styling-base presidium-layouts-base presidium-layouts-blog -x '*.git*' '*/.*' '*/.github/*' 2>/dev/null || true
	@echo "Themes bundle created: themes.zip"

build: ## Build the presidium binary
	@$(MAKE) prepare-themes
	@go build -tags extended -o $(FILENAME) .

test: ## Run tests
	@mkdir -p reports
	@$(MAKE) prepare-themes
	@go test -race -timeout 120s ./...

test-offline: ## Test offline build in Docker container (requires Docker)
	@./test-offline-build.sh

serve-offline: ## Run presidium server in Docker with embedded themes — browse at http://localhost:3131
	@$(MAKE) prepare-themes
	@docker rm -f presidium-offline-serve 2>/dev/null || true
	@DOCKER_BUILDKIT=1 docker build -f Dockerfile.offline-serve -t presidium-offline-serve .
	@echo ""
	@echo "  Presidium server starting at http://localhost:3131"
	@echo "  Press Ctrl+C to stop"
	@echo ""
	@docker run --rm -p 3131:3131 --name presidium-offline-serve presidium-offline-serve

fmt: ## Format Go source files
	go fmt ./...

vet: ## Run go vet
	@$(MAKE) prepare-themes
	@go vet ./...

tidy: ## Tidy and verify module dependencies
	go mod tidy && go mod verify

clean: ## Remove build artifacts
	rm -fr "dist" "$(FILENAME)" "presidium-test" themes.zip

coverage_report: ## Open coverage report in browser
	@go test -coverprofile=reports/tests-cov.out ./... && go tool cover -html=reports/tests-cov.out -o reports/coverage.html && open reports/coverage.html

dist: ## Build distribution binary
	@$(MAKE) prepare-themes
	@mkdir -p "dist" && go build -trimpath -o "dist/presidium" -tags extended

serve-docs: build ## Serve the documentation site
	cd $(DOCSDIR) && make serve

lint: ## Run golangci-lint
	@$(MAKE) prepare-themes
	@golangci-lint run --timeout 10m