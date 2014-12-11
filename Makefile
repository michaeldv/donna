build:
	go build -gcflags -B ./cmd/donna.go

install:
	go install -gcflags -B ./cmd/donna.go

run:
	go run -gcflags -B ./cmd/donna.go -i

test:
	go test
