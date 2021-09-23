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

dist:
	[ -d "dist" ] || mkdir "dist"
	go build -o "dist/presidium" --tags extended
