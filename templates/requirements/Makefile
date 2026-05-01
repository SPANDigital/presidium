default:
	make serve

hugo:
	hugo mod get
	hugo --templateMetrics --ignoreCache --logLevel info

drafts:
	hugo mod get
	hugo --templateMetrics --ignoreCache --logLevel info --buildDrafts

tidy:
	go mod tidy
	hugo mod tidy

refresh: 
	rm -rf public resources go.sum .hugo_build.lock themes
	hugo mod clean
	make tidy
	make hugo

serve:
	make refresh
	hugo server -w --ignoreCache --disableFastRender --logLevel info

serve-drafts:
	make refresh
	hugo server -w --ignoreCache --disableFastRender --logLevel info --buildDrafts

serve-a:
	make refresh
	hugo server -w --ignoreCache --disableFastRender --logLevel info -p 6060

serve-b:
	make refresh
	hugo server -w --ignoreCache --disableFastRender --logLevel info -p 7070

stats:
	python3 scripts/snapshot-stats.py

build:
	make refresh
	hugo build --minify --noBuildLock
	$(MAKE) copy-features
	$(MAKE) copy-users

###############################################################################
# After `hugo build`, copy relevant generated `.feature` files `public/` -> `tests/features`
###############################################################################

SOURCE_DIR=public
TARGET_DIR=src/tests

copy-features:
	@echo "Flattening feature files tree in $(TARGET_DIR)/features..."
	@rm -rf $(TARGET_DIR)/features
	@find $(SOURCE_DIR)/features -type f -name '*.feature' | while read f; do \
		rel=$${f#$(SOURCE_DIR)/features/}; rel=$${rel#*/}; \
		if [[ "$$rel" == */index.feature ]]; then \
			first_line=$$(head -n 1 "$$f" | tr -d '\r'); \
			if echo "$$first_line" | grep -Eq '^(Feature:|@user:|@scenario)'; then \
				dest=$(TARGET_DIR)/features/$${rel%/index.feature}.feature; \
			else \
				continue; \
			fi; \
		else \
			dest=$(TARGET_DIR)/features/$$rel; \
		fi; \
		mkdir -p "$$(dirname "$$dest")"; \
		cp "$$f" "$$dest"; \
	done
	@echo "Copied feature files to $(TARGET_DIR)/features"

copy-users:
	@mkdir -p $(TARGET_DIR)
	@cp $(SOURCE_DIR)/users.json $(TARGET_DIR)/users.json
	@echo "Copied users.json to $(TARGET_DIR)/"