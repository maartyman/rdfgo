f ?= ./...
n ?= .
cpu ?= -1

# Define a list of files to ignore
IGNORE_FILES = cmd/ yaccpar nquads.y

VERSION := $(shell git describe --tags --abbrev=0) # Get the latest tag (e.g., v1.0.0)
MAJOR := $(shell echo $(VERSION) | awk -F'[v.]' '{print $$2}')
MINOR := $(shell echo $(VERSION) | awk -F'[v.]' '{print $$3}')
PATCH := $(shell echo $(VERSION) | awk -F'[v.]' '{print $$4}')

test:
	# Run tests
	@if echo "$(f)" | grep -qE '\.go$$'; then \
		echo "Searching for test file: $(f)"; \
		pkgs=$$(find . -name "$(f)" | xargs -n1 dirname | sort -u); \
		for pkg in $$pkgs; do \
			go test $$pkg -covermode=atomic -coverprofile=covprofile.tmp || exit 1; \
		done; \
	else \
		go test $(f) -covermode=atomic -coverprofile=covprofile.tmp || exit 1; \
	fi; \
	mv covprofile.tmp covprofile
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
	@if echo "$(f)" | grep -qE '\.go$$'; then \
		echo "Searching for test file: $(f)"; \
		pkgs=$$(find . -name "$(f)" | xargs -n1 dirname | sort -u); \
		for pkg in $$pkgs; do \
			go test $$pkg -v -covermode=atomic -coverprofile=covprofile.tmp || exit 1; \
		done; \
	else \
		go test $(f) -v -covermode=atomic -coverprofile=covprofile.tmp || exit 1; \
	fi; \
	mv covprofile.tmp covprofile
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

test-cover:
	# Run tests even if not 100% coverage
	-@$(MAKE) --no-print-directory test || true
	# Generate coverage report
	@go tool cover -html=covprofile -o coverage.html

test-race:
	# Run tests with race detector
	@if echo "$(f)" | grep -qE '\.go$$'; then \
		echo "Searching for test file: $(f)"; \
		pkgs=$$(find . -name "$(f)" | xargs -n1 dirname | sort -u); \
		for pkg in $$pkgs; do \
			go test $$pkg -race -covermode=atomic -coverprofile=covprofile.tmp || exit 1; \
		done; \
	else \
		go test $(f) -race -covermode=atomic -coverprofile=covprofile.tmp || exit 1; \
	fi; \
	mv covprofile.tmp covprofile
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


benchmark:
	# Benchmark
	@if [ "$(f)" = "./..." ]; then \
		f=./performance; \
	else \
		f=*/$(f); \
		echo "Running benchmark in $$f"; \
	fi; \
	if [ "$(n)" = "." ]; then \
		n=.; \
	else \
		n=$(n)$$; \
		echo "Running benchmark $$n"; \
	fi; \
	if [ "$(cpu)" = "-1" ]; then \
		go test $$f -bench=$$n; \
	else \
		go test $$f -bench=$$n -cpu=$(cpu); \
	fi

benchmark-mem:
	# Benchmark
	@if [ "$(f)" = "./..." ]; then \
		f=./performance; \
	else \
		f=*/$(f); \
		echo "Running benchmark in $$f"; \
	fi; \
	if [ "$(n)" = "." ]; then \
		n=.; \
	else \
		n=$(n)$$; \
		echo "Running benchmark $$n"; \
	fi; \
	if [ "$(cpu)" = "-1" ]; then \
		go test $$f -bench=$$n -benchmem; \
	else \
		go test $$f -bench=$$n -benchmem -cpu=$(cpu); \
	fi

flamegraph-cpu:
	# Generate flame graph
	@if [ "$(f)" = "./..." ]; then \
		f=./performance; \
	else \
		f=*/$(f); \
		echo "Running benchmark in $$f"; \
	fi; \
	if [ "$(n)" = "." ]; then \
		n=.; \
	else \
		n=$(n)$$; \
		echo "Running benchmark $$n"; \
	fi; \
	if [ "$(cpu)" = "-1" ]; then \
		go test $$f -bench=$$n -cpuprofile cpu.prof; \
	else \
		go test $$f -bench=$$n -cpu=$(cpu) -cpuprofile cpu.prof; \
	fi
	@go tool pprof -http=:8080 cpu.prof

flamegraph-mem:
	# Generate flame graph
	@if [ "$(f)" = "./..." ]; then \
		f=./performance; \
	else \
		f=*/$(f); \
		echo "Running benchmark in $$f"; \
	fi; \
	if [ "$(n)" = "." ]; then \
		n=.; \
	else \
		n=$(n)$$; \
		echo "Running benchmark $$n"; \
	fi; \
	if [ "$(cpu)" = "-1" ]; then \
		go test $$f -bench=$$n -benchmem -memprofile mem.prof; \
	else \
		go test $$f -bench=$$n -benchmem -cpu=$(cpu) -memprofile mem.prof;  \
	fi
	@go tool pprof -http=:8080 mem.prof

shorten:
	# Shorten lines
	@if [ "$(f)" = "./..." ]; then \
		f=./..; \
	else \
		f=*/$(f); \
		echo "Running benchmark on $$f"; \
	fi; \
	golines $$f -w -m 120

lint:
	# Run linter
	@golangci-lint run

lint-fix:
	# Run linter
	@golangci-lint run --fix

fmt:
	# Format code
	@gofmt -s -w .

pre-commit: build-parser fmt lint test-race

setup-for-release:
	@git checkout master
	@git fetch
	@git pull

bump-version-patch: setup-for-release
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

bump-version-minor: setup-for-release
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

bump-version-major: setup-for-release
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
	@go install github.com/git-chglog/git-chglog/cmd/git-chglog@v0.15.4
	@go install golang.org/x/tools/cmd/goyacc@v0.33.0
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8

build-parser:
	@cd ./lib/parser/nquads && goyacc -o ./yacc.go ./nquads.y
	@cd ./lib/parser/turtle && goyacc -o ./yacc.go ./turtle.y
