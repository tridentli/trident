ifndef $GIT_BRANCH
	GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || printf "nogit")
endif

all:
	@echo "All does nothing"

help:
	@echo "check            - runs: go vet/fmt"
	@echo "tests            - Runs all Golang based tests"
	@echo "vtests           - Runs all Golang based tests (verbose)"
	@echo "versions         - Show pre-selected versions"

versions:
	@echo "Trident: $(GIT_BRANCH)"
	@echo
	@printf "Go version: "; go version
	@printf "Trident branch: "; git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "nogit"; printf "Trident Git commit: "; git rev-parse HEAD 2>/dev/null || echo "nogit"
	@echo

check:
	@echo "Running 'go vet'"
	@go vet ./...
	@echo "Running 'go fmt'"
	@go fmt ./...

tests:
	@echo "Running 'go test'..."
	@go test ./...

vtests:
	@echo "Running 'go test' (verbose)..."
	@go test -v ./...

.PHONY : all help check tests vtests versions
