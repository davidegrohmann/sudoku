GIT_VERSION := $(shell git describe --abbrev=8 --dirty --always --tags)
BUILD_TIME := $(shell date --iso-8601=seconds)
LDFLAGS += -X 'main.buildVersion=$(GIT_VERSION)' -X 'main.buildTime=$(BUILD_TIME)' -s -w

PACKAGE_PROD_DEPS := $(shell find pkg -name "*.go" | grep -v "_test.go")
PACKAGE_TEST_DEPS := $(shell find pkg -name "*_test.go")

all: format lint build static_analysis test
.PHONY: all

build: cmd/sudoku.bin
.PHONY: build

%.bin: %.go $(PACKAGE_PROD_DEPS)
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -mod=vendor -o $@ $<

format: .format
.PHONY: format

.format: $(PACKAGE_PROD_DEPS) $(PACKAGE_TEST_DEPS)
	gofmt -d -s -e $<
	gofmt -s -l $< | wc -l | grep "0" >/dev/null
	touch $@

lint: .lint
.PHONY: lint

.lint: $(PACKAGE_PROD_DEPS) $(PACKAGE_TEST_DEPS)
	golint -set_exit_status ./pkg/... ./cmd/...
	touch $@

static_analysis: .static_analysis
.PHONY: static_analysis

.static_analysis: $(PACKAGE_PROD_DEPS) $(PACKAGE_TEST_DEPS)
	go vet -mod=vendor ./pkg/... ./cmd/...
	touch $@

test: .test
.PHONY: test

.test: $(PACKAGE_PROD_DEPS) $(PACKAGE_TEST_DEPS)
	go test -mod=vendor ./pkg/...
	touch $@

clean:
	rm -rf .format .lint .static_analysis .test cmd/*.bin
.PHONY: clean