#
# This file is part of the Valkyrja Framework package.
#
# Copyright (c) 2016-present Melech Mizrachi
#
# Released under the MIT License. See LICENSE.md for details.
#

# Root task runner for this repository — the analog of PHP's
# composer.json, TypeScript's package.json, and Java's Gradle scripts. Drive
# every CI tool through these targets; check here first for exact target names.
#
# golangci-lint is pinned in .github/ci/lint/go.mod (an isolated tool module) and
# run via `go tool` so its version never mixes with the framework's own deps.
GOLANGCI_LINT ?= go tool -modfile=.github/ci/lint/go.mod golangci-lint

# valkyrjalint ships the golangci-lint configuration and the copyright header
# tool that every Valkyrja Go repository runs. It is pinned in the same isolated
# tool module as golangci-lint. See github.com/valkyrjaio/ci-golangcilint-go.
VALKYRJALINT ?= go tool -modfile=.github/ci/lint/go.mod valkyrjalint

# The package identifier this repository's headers name. It is read from the
# copyright header config, which is the one file that records it.
#
# Warning: never write the value here. `_create-repo.yml` rewrites the config and
# the header of every source file when it scaffolds a repository, and it does not
# read this file. A value written here therefore keeps saying `Project Template`
# in a repository that is not the template, and `make header-fix` then rewrites
# every correct header to the wrong package.
#
# `:=` rather than `?=`: the value is read once, at parse time. A recursive
# assignment re-runs `sed` at every reference.
#
# `$(shell ...)` reports no error of its own, so a changed quoting style or a
# moved config file yields an empty value. `require-package-identifier` below
# catches that. Warning: keep the check in a target, never in a top-level
# `ifeq`. Make reads a top-level conditional before it chooses a target, so a
# broken config would end every target, including `make help` and `make test`,
# which never read the identifier.
PACKAGE_IDENTIFIER := $(shell sed -n "s/^IDENTIFIER='\\(.*\\)'$$/\\1/p" .github/ci/copyright-header/config 2>/dev/null)

# The coverage floor, as a percentage. 100 is the definition of done; it is a hard
# floor, never lowered to accommodate a gap (cover the code, or leave it out of the
# build). Override only for a local spot check: `make coverage COVERAGE_FLOOR=90`.
COVERAGE_FLOOR ?= 100

# The path fragment that marks a package as test support. A reusable double lives
# in a `fixtures` package, and it carries no test of its own: a double exists to
# serve a test, so a test for a double asserts a test.
#
# Go attributes coverage to the package under test. A fixture that every component
# test exercises therefore still reports 0%, because no test runs inside its own
# package. Dropping the package from the profile is what keeps the floor a
# statement about the framework's own code.
COVERAGE_EXCLUDE ?= /fixtures/

.DEFAULT_GOAL := ci

.PHONY: ci
ci: tidy-check config-check header-check fmt-check lint coverage ## Run the full CI gate

.PHONY: config-write
config-write: ## Write .golangci.yml from the shared configuration
	$(VALKYRJALINT) config > .golangci.yml

# `diff` rather than a flag on the command: it matches `go mod tidy -diff`, and it
# prints what differs instead of only reporting that something does.
#
# Warning: never pipe the generator straight into `diff`. `make` runs a recipe under
# `/bin/sh`, which has no `pipefail`, so the pipeline reports `diff`'s status. A
# generator that fails writes nothing, `diff` then reports the whole file as removed,
# and the recipe blames a stale file for a crash. Write the output to a file first, so
# the generator's own status ends the recipe and its stderr is what the reader sees.
.PHONY: config-check
config-check: ## Fail where .golangci.yml differs from the shared configuration
	@generated=$$(mktemp); trap 'rm -f "$$generated"' EXIT; \
		$(VALKYRJALINT) config > "$$generated" \
		&& { diff -u .golangci.yml "$$generated" \
			|| { echo 'FAIL  .golangci.yml is stale. Run: make config-write'; exit 1; }; }

.PHONY: require-package-identifier
require-package-identifier:
	@test -n '$(PACKAGE_IDENTIFIER)' || { \
		printf 'Could not read IDENTIFIER from %s\n' .github/ci/copyright-header/config >&2; \
		exit 1; \
	}

.PHONY: header-fix
header-fix: require-package-identifier ## Write the copyright header into every Go file that lacks it
	$(VALKYRJALINT) header -package '$(PACKAGE_IDENTIFIER)' -w .

.PHONY: header-check
header-check: require-package-identifier ## Fail where a Go file carries no copyright header, or the wrong one
	@$(VALKYRJALINT) header -package '$(PACKAGE_IDENTIFIER)' . \
		&& echo 'PASS  every Go file carries the copyright header' \
		|| { echo 'FAIL  the files above need the header. Run: make header-fix'; exit 1; }

.PHONY: build
build: ## Compile every package
	go build ./...

.PHONY: test
test: ## Run tests with the race detector
	go test -race ./...

.PHONY: coverage
coverage: ## Run tests with coverage and fail below COVERAGE_FLOOR
	go test -race -covermode=atomic -coverprofile=coverage.out ./...
# The filter runs on the profile rather than on the package list, so every test
# still runs, and the report that CI uploads drops the same packages the floor
# drops.
#
# Warning: never filter with `grep -v`. `grep` exits 1 where it writes no line,
# which ends the recipe on a profile that holds nothing else. `awk` exits 0.
# `index` rather than a regex, because COVERAGE_EXCLUDE holds a slash and `awk`
# reads a regex between slashes of its own.
	@awk -v exclude='$(COVERAGE_EXCLUDE)' 'index($$0, exclude) == 0' coverage.out > coverage.out.filtered \
		&& mv coverage.out.filtered coverage.out
	go tool cover -func=coverage.out
# `go tool cover` only reports and always exits 0, so without this the profile was
# generated and then ignored — a run at 55% passed exactly like one at 100%.
#
# Statement coverage only: Go has no native branch coverage, so an untested branch
# inside a covered statement stays invisible here and is a review concern (see
# architecture/go/AGENTS.md).
#
# A package of nothing but constants has no statements at all, and `go tool cover`
# prints "total: (statements) 0.0%" for that — identical to genuinely covering
# nothing. Counting the per-function rows tells the two apart: no rows means there
# was nothing to cover, which is a pass, not a 0% failure.
	@go tool cover -func=coverage.out | awk -v floor='$(COVERAGE_FLOOR)' ' \
		/^total:/ { total = $$3; next } \
		{ rows++ } \
		END { \
			if (rows == 0) { print "No statements to cover; nothing to enforce."; exit 0 } \
			sub(/%$$/, "", total); \
			if (total + 0 < floor + 0) { \
				printf "FAIL  statement coverage %.1f%% < %.1f%%\n", total, floor; exit 1 \
			} \
			printf "PASS  statement coverage %.1f%% >= %.1f%%\n", total, floor \
		}'

.PHONY: fmt
fmt: ## Format the codebase
	$(GOLANGCI_LINT) fmt

.PHONY: fmt-check
fmt-check: ## Fail if the codebase is not formatted
	$(GOLANGCI_LINT) fmt --diff

.PHONY: lint
lint: ## Run golangci-lint (static analysis, security, architecture, dead code)
	$(GOLANGCI_LINT) run

.PHONY: tidy
tidy: ## Tidy module dependencies
	go mod tidy

.PHONY: tidy-check
tidy-check: ## Fail if go.mod / go.sum are not tidy
	go mod tidy -diff

.PHONY: verify
verify: ## Verify dependencies against go.sum
	go mod verify

.PHONY: tools
tools: ## Download the pinned CI tool dependencies
	cd .github/ci/lint && go mod download

.PHONY: outdated
outdated: ## List direct dependencies with newer versions available
	go list -u -m -f '{{if and .Update (not .Indirect)}}{{.Path}}: {{.Version}} -> {{.Update.Version}}{{end}}' all

.PHONY: update
update: ## Update all dependencies to latest and tidy
	go get -u ./...
	go mod tidy

.PHONY: help
help: ## List available targets
	@grep -hE '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | \
		awk 'BEGIN{FS=":.*?## "}{printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'
