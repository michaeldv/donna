# Copyright (c) 2014-2018 by Michael Dvorkin. All Rights Reserved.
# Use of this source code is governed by a MIT-style license that can
# be found in the LICENSE file.
#
# I am making my contributions/submissions to this project solely in my
# personal capacity and am not conveying any rights to any intellectual
# property of any third parties.

VERSION = 4.1
GOFLAGS = -gcflags -B
PACKAGE = github.com/michaeldv/donna/cmd/donna

build:
	go build -x -a -o ./bin/donna $(GOFLAGS) $(PACKAGE)

install:
	go install -x $(GOFLAGS) $(PACKAGE)

run:
	go run $(GOFLAGS) ./cmd/donna/main.go -i

test:
	go test

buildall:
	GOOS=darwin  GOARCH=amd64 go build -a $(GOFLAGS) -o ./bin/donna-$(VERSION)-osx-64         $(PACKAGE)
	GOOS=freebsd GOARCH=amd64 go build -a $(GOFLAGS) -o ./bin/donna-$(VERSION)-freebsd-64     $(PACKAGE)
	GOOS=linux   GOARCH=amd64 go build -a $(GOFLAGS) -o ./bin/donna-$(VERSION)-linux-64       $(PACKAGE)
	GOOS=windows GOARCH=amd64 go build -a $(GOFLAGS) -o ./bin/donna-$(VERSION)-windows-64.exe $(PACKAGE)
	GOOS=windows GOARCH=386   go build -a $(GOFLAGS) -o ./bin/donna-$(VERSION)-windows-32.exe $(PACKAGE)
