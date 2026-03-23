BINARY  := cfddns
MODULE  := github.com/pierresantana/cfddns
MAIN    := ./cmd/cfddns
DISTDIR := dist

.PHONY: build run clean test lint build-linux build-darwin build-all install

build: build-all

run: build
	./$(BINARY)

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(DISTDIR)/$(BINARY)-linux-amd64 $(MAIN)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o $(DISTDIR)/$(BINARY)-linux-arm64 $(MAIN)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 go build -o $(DISTDIR)/$(BINARY)-linux-armv5 $(MAIN)

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

install:
	install -Dm755 $(DISTDIR)/$(BINARY)-linux-amd64 /usr/local/bin/$(BINARY)
	install -Dm600 .env /etc/cfddns/env
	install -Dm644 systemd/cfddns.service /etc/systemd/system/cfddns.service
	install -Dm644 systemd/cfddns.timer /etc/systemd/system/cfddns.timer
	systemctl daemon-reload
	systemctl enable --now cfddns.timer
