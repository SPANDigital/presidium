FILENAME=presidium
DOCSDIR=docs
.DEFAULT_GOAL=help
.PHONY: build test dist clean fmt vet tidy coverage_report help update-themes prepare-themes restore-themes lint checks serve-docs

help: ## Display available targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-18s %s\n", $$1, $$2}'

update-themes: ## Update theme submodules to their latest versions
	@echo "Updating theme submodules to latest..."
	@git submodule update --remote themes/presidium-styling-base themes/presidium-layouts-base themes/presidium-layouts-blog
	@echo "Theme submodules updated."

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
	go build -tags extended -o $(FILENAME) . ; status=$$? ; $(MAKE) restore-themes ; rstatus=$$? ; if [ $$status -eq 0 ]; then status=$$rstatus; fi ; exit $$status

test: prepare-themes ## Run tests with coverage
	@mkdir -p reports
	go test -race -timeout 120s ./... -coverprofile=reports/tests-cov.out ; status=$$? ; $(MAKE) restore-themes ; rstatus=$$? ; if [ $$status -eq 0 ]; then status=$$rstatus; fi ; exit $$status

fmt: ## Format Go source files
	go fmt ./...

vet: prepare-themes ## Run go vet
	go vet ./... ; status=$$? ; $(MAKE) restore-themes ; rstatus=$$? ; if [ $$status -eq 0 ]; then status=$$rstatus; fi ; exit $$status

tidy: ## Tidy and verify module dependencies
	go mod tidy && go mod verify

clean: restore-themes ## Remove build artifacts and restore themes
	rm -fr "dist" "$(FILENAME)" "presidium-test"

coverage_report: ## Open coverage report in browser
	@go tool cover -html=reports/tests-cov.out

dist: prepare-themes ## Build distribution binary
	mkdir -p "dist" && go build -trimpath -o "dist/presidium" --tags extended ; status=$$? ; $(MAKE) restore-themes ; rstatus=$$? ; if [ $$status -eq 0 ]; then status=$$rstatus; fi ; exit $$status

serve-docs:
	cd $(DOCSDIR) && make serve

lint: prepare-themes ## Run golangci-lint
	golangci-lint run --timeout 10m ; status=$$? ; $(MAKE) restore-themes ; rstatus=$$? ; if [ $$status -eq 0 ]; then status=$$rstatus; fi ; exit $$status