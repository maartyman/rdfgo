file ?= ./..
test-file ?= ./...

# Define a list of files to ignore
IGNORE_FILES = cmd/

VERSION := $(shell git describe --tags --abbrev=0) # Get the latest tag (e.g., v1.0.0)
MAJOR := $(shell echo $(VERSION) | awk -F'[v.]' '{print $$2}')
MINOR := $(shell echo $(VERSION) | awk -F'[v.]' '{print $$3}')
PATCH := $(shell echo $(VERSION) | awk -F'[v.]' '{print $$4}')

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
	$(eval NEW_VERSION := v$(MAJOR).$(MINOR).$(shell echo $(PATCH) + 1 | bc))
	@echo "Are you sure you want to bump the version to $(NEW_VERSION)? (y/n)" && read ans && [ $${ans:-n} = y ]

	# Generate tag and changelog
	@git reset
	@git-chglog -o CHANGELOG.md
	@git add CHANGELOG.md
	@git commit -am "chore(release): bump to $(NEW_VERSION)"
	@git tag -a $(NEW_VERSION) -m "$(NEW_VERSION)" -m "See https://github.com/maartyman/rdfgo/blob/$(NEW_VERSION)/CHANGELOG.md for changes."
	@git-chglog -o CHANGELOG.md
	@git add CHANGELOG.md
	@git commit --amend --no-edit
	@git push origin HEAD
	@git push origin $(NEW_VERSION)

	# Create GitHub release
	@gh release create $(NEW_VERSION) --title "Release $(NEW_VERSION)" --notes "$$(cat CHANGELOG.md)"

bump-version-minor:
	# Bump minor version
	$(eval NEW_VERSION := v$(MAJOR).$(shell echo $(MINOR) + 1 | bc).0)
	@echo "Are you sure you want to bump the version to $(NEW_VERSION)? (y/n)" && read ans && [ $${ans:-n} = y ]

	# Generate tag and changelog
	@git reset
	@git-chglog -o CHANGELOG.md
	@git add CHANGELOG.md
	@git commit -am "chore(release): bump to $(NEW_VERSION)"
	@git tag -a $(NEW_VERSION) -m "$(NEW_VERSION)" -m "See https://github.com/maartyman/rdfgo/blob/$(NEW_VERSION)/CHANGELOG.md for changes."
	@git-chglog -o CHANGELOG.md
	@git add CHANGELOG.md
	@git commit --amend --no-edit
	@git push origin HEAD
	@git push origin $(NEW_VERSION)

	# Create GitHub release
	@gh release create $(NEW_VERSION) --title "Release $(NEW_VERSION)" --notes "$$(cat CHANGELOG.md)"

bump-version-major:
	# Bump major version and update go.mod
	$(eval NEW_VERSION := v$(shell echo $(MAJOR) + 1 | bc).0.0)
	@echo "Are you sure you want to bump the version to $(NEW_VERSION)? (y/n)" && read ans && [ $${ans:-n} = y ]

	# Update go.mod for new major version
	@git reset
	@sed -i'' -e 's/^module \(.*\)/module \1\/v$(shell echo $(MAJOR) + 1 | bc)/' go.mod
	@find . -name '*.go' -type f -exec sed -i'' -e 's/\(.*\)\/v$(MAJOR)/\1\/v$(shell echo $(MAJOR) + 1 | bc)/g' {} \;
	@git add .

	# Generate tag and changelog
	@git-chglog -o CHANGELOG.md
	@git add CHANGELOG.md
	@git commit -am "chore(release): bump to $(NEW_VERSION) and update go.mod for v$(shell echo $(MAJOR) + 1 | bc)"
	@git push origin HEAD
	@git tag -a $(NEW_VERSION) -m "$(NEW_VERSION)" -m "See https://github.com/maartyman/rdfgo/blob/$(NEW_VERSION)/CHANGELOG.md for changes."
	@git-chglog -o CHANGELOG.md
	@git add CHANGELOG.md
	@git commit --amend --no-edit
	@git push origin HEAD
	@git push origin $(NEW_VERSION)

	# Create GitHub release
	@gh release create $(NEW_VERSION) --title "Release $(NEW_VERSION)" --notes "$$(cat CHANGELOG.md)"

setup-project:
	# Make all files in .githooks executable
	@chmod +x .githooks/*
	# Setup git hooks
	@git config core.hooksPath .githooks
	# Install Dependencies
	@go install github.com/git-chglog/git-chglog/cmd/git-chglog@latest
