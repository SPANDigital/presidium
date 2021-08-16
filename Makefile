FILENAME=main

test:
	go test ./...

build:
	go get -u github.com/gobuffalo/packr/v2/packr2
	packr2
	go build  -o $(FILENAME) main.go
	packr2 clean

default:
	mkdir -p "dist"
	go build -o "dist/presidium"
