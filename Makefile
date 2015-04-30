version = 2.0rc1
goflags = -gcflags -B

build:
	go build $(goflags) ./cmd/donna.go

install:
	go install $(goflags) ./cmd/donna.go

run:
	go run $(goflags) ./cmd/donna.go -i

test:
	go test

buildall:
	GOOS=darwin  GOARCH=amd64 go build $(goflags) -o donna-$(version)-osx-64     ./cmd/donna.go
	GOOS=freebsd GOARCH=amd64 go build $(goflags) -o donna-$(version)-freebsd-64 ./cmd/donna.go
	GOOS=linux   GOARCH=amd64 go build $(goflags) -o donna-$(version)-linux-64   ./cmd/donna.go
	GOOS=windows GOARCH=amd64 go build $(goflags) -o donna-$(version)-windows-64 ./cmd/donna.go
