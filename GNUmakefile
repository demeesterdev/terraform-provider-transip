TEST?=$$(go list ./... |grep -v 'vendor'|grep -v 'examples')
PKG_NAME=transip

#make sure we catch schema errors during testing
TF_SCHEMA_PANIC_ON_ERROR=1
GO111MODULE=on
GOFLAGS=-mod=vendor

default: build

build: fmtcheck
	go install

test: fmtcheck
	go test -mod=vendor -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 --format short-verbose -- $(TESTARGS) -timeout=30s -parallel=4 -mod=vendor -coverprofile=cover.out

testcover: test
	go tool cover -html=cover.out


testacc: fmtcheck
	TF_ACC=1 gotestsum --format short-verbose -- $(TEST) $(TESTARGS) -timeout 180m -coverprofile=cover-acc.out

testacccover: testacc
	go tool cover -html=cover-acc.out

# Currently required by tf-deploy compile
fmtcheck:
	@sh "$(CURDIR)/scripts/gofmtcheck.sh"

fmt:
	@echo "==> Fixing source code with gofmt..."
	# This logic should match the search logic in scripts/gofmtcheck.sh
	find . -name '*.go' | grep -v vendor | xargs gofmt -s -w

lint:
	@echo "==> Checking source code against linters..."
	golangci-lint run ./...

tools:
	@echo "==> installing required tooling..."
	@sh "$(CURDIR)/scripts/gogetcookie.sh"
	GO111MODULE=off go get -u github.com/client9/misspell/cmd/misspell
	GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	GO111MODULE=off go get -u gotest.tools/gotestsum

vendor:
	go mod tidy
	go mod download
	go mod vendor

.PHONY: build test testcover testacc testacccover fmtcheck fmt lint tools vendor
