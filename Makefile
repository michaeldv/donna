VERSION = 3.1
GOFLAGS = -gcflags -B

build:
	go build $(GOFLAGS) ./cmd/donna.go

install:
	go install $(GOFLAGS) ./cmd/donna.go

run:
	go run $(GOFLAGS) ./cmd/donna.go -i

test:
	go test

buildall:
	GOOS=darwin  GOARCH=amd64 go build $(GOFLAGS) -o donna-$(VERSION)-osx-64         ./cmd/donna.go
	GOOS=freebsd GOARCH=amd64 go build $(GOFLAGS) -o donna-$(VERSION)-freebsd-64     ./cmd/donna.go
	GOOS=linux   GOARCH=amd64 go build $(GOFLAGS) -o donna-$(VERSION)-linux-64       ./cmd/donna.go
	GOOS=windows GOARCH=amd64 go build $(GOFLAGS) -o donna-$(VERSION)-windows-64.exe ./cmd/donna.go
	GOOS=windows GOARCH=386   go build $(GOFLAGS) -o donna-$(VERSION)-windows-32.exe ./cmd/donna.go
