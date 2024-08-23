file ?= ./..
test-file ?= ./...

test:
	go test $(test-file)

test-verbose:
	go test $(test-file) -v

test-cover:
	go test $(test-file) -coverprofile=c.out && go tool cover -html=c.out

shorten:
	golines $(file) -w -m 120

lint:
	golangci-lint run

fmt:
	gofmt -s -w .

pre-commit: fmt lint test
