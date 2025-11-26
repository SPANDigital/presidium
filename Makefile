FILENAME=presidium
TESTDIRS=`go list ./... |grep -v "vendor/" |grep -v "swagger/"`
.DEFAULT_GOAL=build
.PHONY: dist clean

test:
	@mkdir -p reports
	go test -p 1 -v $(TESTDIRS) -coverprofile=reports/tests-cov.out


pack:
	go mod tidy

build:
	make pack
	go build -tags extended -o $(FILENAME) main.go


clean:
	rm -fr "dist"

coverage_report:
	@go tool cover -html=reports/tests-cov.out

dist:
	[ -d "dist" ] || mkdir "dist"
	go build -o "dist/presidium" --tags extended
