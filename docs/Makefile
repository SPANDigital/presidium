PRESIDIUM=../presidium

default:
	make serve

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf public resources go.sum .hugo_build.lock themes
	@echo "Build artifacts cleaned."

hugo:
	$(PRESIDIUM) hugo mod get
	$(PRESIDIUM) hugo --templateMetrics --ignoreCache --logLevel info

drafts:
	$(PRESIDIUM) hugo mod get
	$(PRESIDIUM) hugo --templateMetrics --ignoreCache --logLevel info --buildDrafts

tidy:
	go mod tidy
	$(PRESIDIUM) hugo mod tidy

refresh:
	$(PRESIDIUM) refresh

serve:
	$(PRESIDIUM) refresh
	$(PRESIDIUM) hugo server -w --ignoreCache --disableFastRender --logLevel info

serve-proxy:
	$(PRESIDIUM) refresh
	$(PRESIDIUM) server -- -w --ignoreCache --disableFastRender --logLevel info

serve-drafts:
	$(PRESIDIUM) refresh
	$(PRESIDIUM) hugo server -w --ignoreCache --disableFastRender --logLevel info --buildDrafts

serve-drafts-proxy:
	$(PRESIDIUM) refresh
	$(PRESIDIUM) server -- -w --ignoreCache --disableFastRender --logLevel info --buildDrafts

serve-a:
	$(PRESIDIUM) refresh
	$(PRESIDIUM) hugo server -w --ignoreCache --disableFastRender --logLevel info -p 6060

serve-b:
	$(PRESIDIUM) refresh
	$(PRESIDIUM) hugo server -w --ignoreCache --disableFastRender --logLevel info -p 7070