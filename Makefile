GOCMD=go
GOCLEAN=$(GOCMD) clean
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOFMT=$(GOCMD) fmt
GOTEST=$(GOCMD) test
CLI_BINARY_NAME=faas-cli

cli:
	$(GOFMT) ./...
	$(GOBUILD) -v -o ./faas-cli/$(CLI_BINARY_NAME) ./faas-cli
install:
	# faas-cli
	$(GOFMT) ./...
	$(GOINSTALL) -v ./faas-cli
clean:
	rm -rf faas-cli/$(CLI_BINARY_NAME)
