default:
	make hugo

clear:
	rm -rf public resources

hugo:
	make clear
	hugo mod get
	hugo --templateMetrics --ignoreCache --logLevel info

drafts:
	make clear
	hugo mod get
	hugo --templateMetrics --ignoreCache --logLevel info --buildDrafts

tidy:
	make clear
	go mod tidy
	hugo mod tidy

refresh:
	make clear
	hugo mod clean
	make tidy
	make hugo

serve:
	make clear
	hugo server -w --ignoreCache --disableFastRender --logLevel info

serve-a:
	make clear
	hugo server -w --ignoreCache --disableFastRender --logLevel info -p 6060

serve-b:
	make clear
	hugo server -w --ignoreCache --disableFastRender --logLevel info -p 7070