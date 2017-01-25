TRIDENT := ${PWD}
GOPATH := ${PWD}/ext/_gopath
export GOPATH

ifndef $GIT_BRANCH
	GIT_BRANCH := `git rev-parse --abbrev-ref HEAD 2>/dev/null || printf "nogit"`
endif

# Pull correct branches for PF from local files.
PF_GIT_BRANCH = issue_67
PF_GIT_REPO = https://github.com/tridentli/pitchfork.git

# dpkg-buildpackage calls make, so <all> should be empty.
all:
	@echo "All does nothing"

help:
	@echo "all              - Build it all (called from dpkg-buildpackage)"
	@echo "clean_ext        - cleans 'ext' directory"
	@echo "pkg              - ext+deps+pkg_only"
	@echo "pkg_only         - Only builds Debian package (no dependency updates)"
	@echo "versions         - Show pre-selected versions"
	@echo "ext              - Fetches initial external dependencies (for use by local devs or trident-ext-pkg)"
	@echo "deps             - Updates external dependencies (for use by local devs or trident-ext-pkg)"
	@echo "build_ext        - Create a GOPATH version of the local trident dir"
	@echo "check            - runs: go vet/fmt, also on Pitchfork"
	@echo "tests		- Runs all Golang based tests"
	@echo "vtests		- Runs all Golang based tests (verbose)"

clean_ext:
	@echo "Cleansing 'ext'"
	@rm -rf ext

pkg: ext deps pkg_only

pkg_only: build_ext
	@echo "Starting package build..."
	dpkg-buildpackage -uc -us -F
	@echo "Starting package build - done"

versions:
	@echo "Trident: $(GIT_BRANCH)"
	@echo "Pitchfork: $(PF_GIT_BRANCH)"
	@echo
	@printf "Go version: "; go version
	@printf "Trident branch: "; git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "nogit"; printf "Trident Git commit: "; git rev-parse HEAD 2>/dev/null || echo "nogit"
	@echo
	@printf "Build OS: "
	@if [ -f /etc/debian_version ]; then printf "Debian "; cat /etc/debian_version; else echo "Not Debian"; fi
	@echo

build_ext:
	@echo "Creating temporary 'ext' for build..."

	@echo
	@echo "Creating 'ext' directory for external dependencies and gopath"
	@mkdir -p ext/_gopath/src/trident.li
	@echo "Symlink Trident into GOPATH"
	@ln -s ../../../../ ext/_gopath/src/trident.li/trident

ext:
	@$(MAKE) versions
	@echo "Retrieving 'ext' dependencies..."

	@echo "Git Clone Pitchfork into GOPATH [ $(PF_GIT_BRANCH) ]"
	@mkdir -p ext/_gopath/src/trident.li
	@git clone -b $(PF_GIT_BRANCH) $(PF_GIT_REPO) ext/_gopath/src/trident.li/pitchfork

	@echo
	@echo "Clone Trident Go Dependencies"
	@mkdir -p ext/_gopath/src/trident.li
	@git clone https://github.com/tridentli/go ext/_gopath/src/trident.li/go

	@echo
	@echo "Retrieving 'ext' dependencies - done"
	@echo

deps: ext
	@$(MAKE) versions
	@echo "Updating 'ext' Dependencies..."

	@echo
	@echo "Updating Pitchfork..."
	@cd ext/_gopath/src/trident.li/pitchfork; (git pull 2>/dev/null && git checkout ${PF_GIT_BRANCH} && printf "Pitchfork branch: " && git rev-parse --abbrev-ref HEAD && printf "Pitchfork GIT commit: " && git rev-parse HEAD) || echo "NO GIT"

	@echo
	@echo "Updating Pitchfork dependencies"
	@cd ext/_gopath/src/trident.li/pitchfork && $(MAKE) deps

	@echo
	@echo "Updating Trident Go Dependencies"
	@cd ext/_gopath/src/trident.li/go; git pull 2>/dev/null || echo "NO GIT"

	@echo
	@echo "Updating Golang dependencies (imports)..."
	@GOPATH=${PWD}/ext/_gopath go get -v -d -t ./...

	@echo
	@echo "Updating Dependencies - done"
	@echo

check:
	@echo "Running 'go vet'"
	@go vet ./...
	@echo "Running 'go fmt'"
	@go fmt ./...

tests:
	@echo "Running 'go test'..."
	@go test ./...
	@cd ext/_gopath/src/trident.li/pitchfork && make tests

vtests:
	@echo "Running 'go test' (verbose)..."
	@go test -v ./...
	@cd ext/_gopath/src/trident.li/pitchfork && make vtests

.PHONY : all help clean_ext pkg pkg_only versions deps check tests vtests
