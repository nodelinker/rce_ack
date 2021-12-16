

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get


build: build-linux

clean: clean-windows clean-linux
	$(GOCLEAN)

clean-windows:
	rm -rf ./build/rce_server.exe

clean-linux:
	rm -rf ./build/rce_server


# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o ./build/rce_server ./cmd/


build-linux-arm:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GOBUILD) -o ./build/rce_server ./cmd/


build-windows:
	$(GOBUILD) -o ./build/rce_server.exe ./cmd/agent_service/


.PHONY: build