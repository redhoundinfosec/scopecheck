BINARY=scopecheck
VERSION=0.1.0
LDFLAGS=-ldflags "-s -w"

.PHONY: build test clean release lint

build:
	go build $(LDFLAGS) -o $(BINARY) ./cmd/scopecheck/

test:
	go test ./... -v -count=1

lint:
	go vet ./...
	gofmt -l .

clean:
	rm -f $(BINARY) $(BINARY)-*

release: clean
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY)-linux-amd64 ./cmd/scopecheck/
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BINARY)-linux-arm64 ./cmd/scopecheck/
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY)-windows-amd64.exe ./cmd/scopecheck/
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY)-darwin-amd64 ./cmd/scopecheck/
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BINARY)-darwin-arm64 ./cmd/scopecheck/
	@echo "Release binaries built for v$(VERSION)"
	@ls -lh $(BINARY)-*
