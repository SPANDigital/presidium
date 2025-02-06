FILENAME=main
TESTDIRS=`go list ./... |grep -v "vendor/" |grep -v "swagger/"`
.DEFAULT_GOAL=build
.PHONY: dist clean

test:
	@mkdir -p reports
	go test -p 1 -v $(TESTDIRS) -coverprofile=reports/tests-cov.out

pack:
	go get -u github.com/gobuffalo/packr/v2/packr2
	packr2
	go mod tidy

build:
	make pack
	go build  -o $(FILENAME) main.go
	packr2 clean

clean:
	rm -fr "dist"

coverage_report:
	@go tool cover -html=reports/tests-cov.out

dist:
	[ -d "dist" ] || mkdir "dist"
	go build -o "dist/presidium" --tags extended

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
	hugo mod clean
	make tidy
	make hugo

serve:
	hugo server -w --ignoreCache --disableFastRender --logLevel info

serve-a:
	hugo server -w --ignoreCache --disableFastRender --logLevel info -p 6060

serve-b:
	hugo server -w --ignoreCache --disableFastRender --logLevel info -p 7070