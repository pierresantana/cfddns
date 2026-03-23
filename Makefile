BINARY  := cfddns
MODULE  := github.com/pierresantana/cfddns
MAIN    := ./cmd/cfddns
DISTDIR := dist

.PHONY: build run clean test lint build-linux build-darwin build-all

build: build-all

run: build
	./$(BINARY)

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(DISTDIR)/$(BINARY)-linux-amd64 $(MAIN)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o $(DISTDIR)/$(BINARY)-linux-arm64 $(MAIN)

build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $(DISTDIR)/$(BINARY)-darwin-amd64 $(MAIN)
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o $(DISTDIR)/$(BINARY)-darwin-arm64 $(MAIN)

build-all: build-linux build-darwin

test:
	go test ./...

lint:
	golangci-lint run ./...

clean:
	rm -f $(BINARY)
	rm -rf $(DISTDIR)
