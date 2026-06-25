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

frontmatter:
	rm -f data/presidium/frontmatter-schema.yaml
	hugo mod get
	hugo --ignoreCache --gc --config config.yml,dependencies.config.yml
	mkdir -p data/presidium
	cp public/frontmatter-schema.yaml data/presidium/frontmatter-schema.yaml
	mkdir -p static/schemas
	yq eval -o=json data/presidium/frontmatter-schema.yaml > static/schemas/frontmatter-schema.json