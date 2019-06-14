TEST?=$$(go list ./... |grep -v 'vendor'|grep -v 'examples')
PKG_NAME=transip

#make sure we catch schema errors during testing
TF_SCHEMA_PANIC_ON_ERROR=1
GO111MODULE=on

default: build

build: fmtcheck
	go install

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -v -timeout=30s -parallel=4 

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 180m 

# Currently required by tf-deploy compile
fmtcheck:
	@sh "$(CURDIR)/scripts/gofmtcheck.sh"

lint: tools
	@echo "==> Checking source code against linters..."
	golangci-lint run ./...



tools:
	@echo "==> installing required tooling..."
	go get -u github.com/client9/misspell/cmd/misspell
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint

cache:
	@echo "==> priming cache..."
	go mod download
