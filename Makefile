file ?= ./..
test-file ?= ./...

# Define a list of files to ignore
IGNORE_FILES = cmd/

test:
	# Run tests
	@go test $(test-file) -covermode atomic -coverprofile=covprofile
	@for pattern in $(IGNORE_FILES); do \
		grep -v -E $$pattern covprofile > tmp_filtered.out; \
		mv tmp_filtered.out covprofile; \
	done;
	@coverage=$$(go tool cover -func=covprofile | grep total | awk '{print $$3}' | sed 's/%//'); \
	if [ $${coverage%.*} -ne 100 ]; then \
		echo "Total test coverage is not 100%: $$coverage%"; \
		exit 1; \
	else \
		echo "Total test coverage is: $$coverage%"; \
	fi

test-verbose:
	# Run tests verbose
	@go test $(test-file) -v -covermode atomic -coverprofile=covprofile
	@for pattern in $(IGNORE_FILES); do \
		grep -v -E $$pattern covprofile > tmp_filtered.out; \
		mv tmp_filtered.out covprofile; \
	done;
	@coverage=$$(go tool cover -func=covprofile | grep total | awk '{print $$3}' | sed 's/%//'); \
	if [ $${coverage%.*} -ne 100 ]; then \
		echo "Total test coverage is not 100%: $$coverage%"; \
		exit 1; \
	else \
		echo "Total test coverage is: $$coverage%"; \
	fi

test-cover: test
	# Generate coverage report
	@go tool cover -html=covprofile

test-cover-save: test
	# Generate coverage report
	@go tool cover -html=covprofile -o coverage.html

test-race:
	# Run tests with race detector
	@go test $(test-file) -race -covermode atomic -coverprofile=covprofile
	@for pattern in $(IGNORE_FILES); do \
		grep -v -E $$pattern covprofile > tmp_filtered.out; \
		mv tmp_filtered.out covprofile; \
	done;
	@coverage=$$(go tool cover -func=covprofile | grep total | awk '{print $$3}' | sed 's/%//'); \
	if [ $${coverage%.*} -ne 100 ]; then \
		echo "Total test coverage is not 100%: $$coverage%"; \
		exit 1; \
	else \
		echo "Total test coverage is: $$coverage%"; \
	fi

shorten:
	# Shorten lines
	@golines $(file) -w -m 120

lint:
	# Run linter
	@golangci-lint run

fmt:
	# Format code
	@gofmt -s -w .

pre-commit: fmt lint test-race

setup-project:
	# Make all files in .githooks executable
	@chmod +x .githooks/*
	# Setup git hooks
	@git config core.hooksPath .githooks

