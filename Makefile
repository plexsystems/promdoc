.PHONY: build
build:
	go build

.PHONY: test
test:
	go test -v ./...
