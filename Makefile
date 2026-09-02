.PHONY: fmt vet test build check

fmt:
	gofmt -w $$(find . -name '*.go' -not -path './vendor/*')

vet:
	go vet ./...

test:
	go test ./...

build:
	mkdir -p bin
	go build -trimpath -o bin/vardiel ./cmd/vardiel

check:
	test -z "$$(gofmt -l .)"
	go vet ./...
	go test ./...
	go build ./cmd/vardiel
