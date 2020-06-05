.PHONY: build
build:
	go build

.PHONY: test
test:
	go test -v ./...

.PHONY: release
release:
	GOOS=darwin GOARCH=amd64 go build -o promdoc-darwin-amd64
	GOOS=windows GOARCH=amd64 go build -o promdoc-windows-amd64
	GOOS=linux GOARCH=amd64 go build -o promdoc-linux-amd64
