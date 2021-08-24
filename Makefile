FILENAME=main

.DEFAULT_GOAL=build
.PHONY: dist clean

test:
	go test ./...

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
	go build -o "dist/presidium"
