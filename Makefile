# Copyright (c) 2014-2016 by Michael Dvorkin. All Rights Reserved.
# Use of this source code is governed by a MIT-style license that can
# be found in the LICENSE file.

VERSION = 3.1
GOFLAGS = -gcflags -B
PACKAGE = github.com/michaeldv/donna/cmd/donna

build:
	go build -x $(GOFLAGS) $(PACKAGE)

install:
	go install -x $(GOFLAGS) $(PACKAGE)

run:
	go run $(GOFLAGS) ./cmd/donna/main.go -i

test:
	go test

buildall:
	GOOS=darwin  GOARCH=amd64 go build $(GOFLAGS) -o ./bin/donna-$(VERSION)-osx-64         $(PACKAGE)
	GOOS=freebsd GOARCH=amd64 go build $(GOFLAGS) -o ./bin/donna-$(VERSION)-freebsd-64     $(PACKAGE)
	GOOS=linux   GOARCH=amd64 go build $(GOFLAGS) -o ./bin/donna-$(VERSION)-linux-64       $(PACKAGE)
	GOOS=windows GOARCH=amd64 go build $(GOFLAGS) -o ./bin/donna-$(VERSION)-windows-64.exe $(PACKAGE)
	GOOS=windows GOARCH=386   go build $(GOFLAGS) -o ./bin/donna-$(VERSION)-windows-32.exe $(PACKAGE)
