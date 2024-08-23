file ?= ./..
test-file ?= ./...

# Define a list of files to ignore
IGNORE_FILES = cmd/

test:
	# Run tests
	@go test $(test-file) -coverprofile=c.out
	@for pattern in $(IGNORE_FILES); do \
		grep -v -E $$pattern c.out > tmp_filtered.out; \
		mv tmp_filtered.out c.out; \
	done;
	@coverage=$$(go tool cover -func=c.out | grep total | awk '{print $$3}' | sed 's/%//'); \
	if [ $${coverage%.*} -ne 100 ]; then \
		echo "Total test coverage is not 100%: $$coverage%"; \
		exit 1; \
	else \
		echo "Total test coverage is: $$coverage%"; \
	fi

test-verbose:
	# Run tests verbose
	@go test $(test-file) -v -coverprofile=c.out
	@for pattern in $(IGNORE_FILES); do \
		grep -v -E $$pattern c.out > tmp_filtered.out; \
		mv tmp_filtered.out c.out; \
	done;
	@coverage=$$(go tool cover -func=c.out | grep total | awk '{print $$3}' | sed 's/%//'); \
	if [ $${coverage%.*} -ne 100 ]; then \
		echo "Total test coverage is not 100%: $$coverage%"; \
		exit 1; \
	else \
		echo "Total test coverage is: $$coverage%"; \
	fi

test-cover: test
	# Generate coverage report
	@go tool cover -html=c.out

test-cover-save: test
	# Generate coverage report
	@go tool cover -html=c.out -o coverage.html

shorten:
	# Shorten lines
	@golines $(file) -w -m 120

lint:
	# Run linter
	@golangci-lint run

fmt:
	# Format code
	@gofmt -s -w .

pre-commit: fmt lint test
