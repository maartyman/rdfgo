file ?= ./..
test-file ?= ./...

# Define a list of files to ignore
IGNORE_FILES = cmd/

VERSION := $(shell git describe --tags --abbrev=0) # Get the latest tag (e.g., v1.0.0)
MAJOR := $(shell echo $(VERSION) | cut -d. -f1 | sed 's/v//') # Extract major version
MINOR := $(shell echo $(VERSION) | cut -d. -f2) # Extract minor version
PATCH := $(shell echo $(VERSION) | cut -d. -f3) # Extract patch version

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

bump-version-patch:
	# Bump patch version
	@NEW_VERSION=v$(MAJOR).$(MINOR).$$(($$(echo $(PATCH))+1))
	@git-chglog -o CHANGELOG.md
	@git commit -am "chore(release): bump to $$NEW_VERSION"
	@git tag $$NEW_VERSION
	@git push origin $$NEW_VERSION

bump-version-minor:
	# Bump minor version
	@NEW_VERSION=v$(MAJOR).$$(($$(echo $(MINOR))+1)).0
	@git-chglog -o CHANGELOG.md
	@git commit -am "chore(release): bump to $$NEW_VERSION"
	@git tag $$NEW_VERSION
	@git push origin $$NEW_VERSION

bump-version-major:
	# Bump major version and update go.mod
	@NEW_VERSION=v$$(($$(echo $(MAJOR))+1)).0.0
	@git-chglog -o CHANGELOG.md
	# Update go.mod for new major version
	@sed -i'' -e 's/^module \(.*\)/module \1\/v$$(($$(echo $(MAJOR))+1))/' go.mod
	# Update imports for the new major version
	@find . -name '*.go' -type f -exec sed -i'' -e 's/\(.*\)\/v$(MAJOR)/\1\/v$$(($$(echo $(MAJOR))+1))/g' {} \;
	@git commit -am "chore(release): bump to $$NEW_VERSION and update go.mod for v$$(($$(echo $(MAJOR))+1))"
	@git tag $$NEW_VERSION
	@git push origin $$NEW_VERSION

setup-project:
	# Make all files in .githooks executable
	@chmod +x .githooks/*
	# Setup git hooks
	@git config core.hooksPath .githooks
	# Install Dependencies
	@go install github.com/git-chglog/git-chglog/cmd/git-chglog@latest

