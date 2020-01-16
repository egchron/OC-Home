all: fmt vet install test

fmt:
	go fmt ./...

vet:
	go vet ./...

install:
	go install ./...

test:
	go test ./...

	~/go/bin/server -help
