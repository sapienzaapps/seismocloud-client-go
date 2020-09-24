.PHONY: test

test:
	go test -v ./...
	go vet ./...
	gosec ./...
	staticcheck ./...
	ineffassign .
	errcheck ./...