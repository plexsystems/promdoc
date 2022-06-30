.PHONY: build
build:
	@go build

.PHONY: test
test:
	@go test -v ./... -count=1

.PHONY: release
release:
	@test $(version)
	GOOS=darwin GOARCH=amd64 go build -o promdoc-darwin-amd64 -ldflags="-X 'github.com/plexsystems/promdoc/internal/commands.version=$(version)'"
	GOOS=darwin GOARCH=arm64 go build -o promdoc-darwin-arm64 -ldflags="-X 'github.com/plexsystems/promdoc/internal/commands.version=$(version)'"
	GOOS=windows GOARCH=amd64 go build -o promdoc-windows-amd64 -ldflags="-X 'github.com/plexsystems/promdoc/internal/commands.version=$(version)'"
	GOOS=linux GOARCH=amd64 go build -o promdoc-linux-amd64 -ldflags="-X 'github.com/plexsystems/promdoc/internal/commands.version=$(version)'"
