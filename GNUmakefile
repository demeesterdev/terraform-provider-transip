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
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4 -mod=vendor

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 180m -mod=vendor

# Currently required by tf-deploy compile
fmtcheck:
	@sh "$(CURDIR)/scripts/gofmtcheck.sh"

lint:
	@echo "==> Checking source code against linters..."
	golangci-lint run ./...

tools:
	@echo "==> installing required tooling..."
	@sh "$(CURDIR)/scripts/gogetcookie.sh"
	GO111MODULE=off go get -u github.com/client9/misspell/cmd/misspell
	GO111MODULE=off go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

vendor:
	go mod tidy
	go mod download
	go mod vendor
