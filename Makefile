.DEFAULT_GOAL := build
install:
	go get gopkg.in/alecthomas/gometalinter.v2
	gometalinter.v2 --install
.PHONY: install


build: install
	gometalinter.v2 ./...
	go test -v ./...

	GOOS=linux GOARCH=amd64 go build -o bootstrap ./example
	@zip -9 -r ./handler.zip bootstrap
.PHONY: build

