PACK            := hcloudplus
ORG             := hatemosphere
PROJECT         := github.com/$(ORG)/pulumi-$(PACK)

PROVIDER        := pulumi-resource-$(PACK)
TFGEN           := pulumi-tfgen-$(PACK)
# VERSION resolution order:
#   1. $VERSION in the environment wins.
#   2. Else, most recent git tag matching v*.
#   3. Else, fallback "0.0.1".
# Must always resolve to a non-empty semver — tfbridge panics on empty.
VERSION         ?= $(shell git describe --tags --abbrev=0 --match 'v*' 2>/dev/null | sed 's/^v//')
ifeq ($(VERSION),)
VERSION         := 0.0.1
endif
PROVIDER_PATH   := provider
VERSION_PKG     := $(PROJECT)/$(PROVIDER_PATH)/pkg/version
SCHEMA_PATH     := $(PROVIDER_PATH)/cmd/$(PROVIDER)/schema.json
LDFLAGS         := -X $(VERSION_PKG).Version=$(VERSION)
BIN             := $(CURDIR)/bin

.PHONY: default help all build provider tfgen schema \
        build_sdks build_nodejs build_python build_go build_dotnet build_java \
        install install_provider install_nodejs_sdk \
        fmt lint lint_fix tidy clean

default: help

help:
	@echo "Targets:"
	@echo "  build            Build everything: schema + provider + all SDKs"
	@echo "  provider         Build $(PROVIDER) binary under bin/"
	@echo "  tfgen            Build $(TFGEN) binary under bin/"
	@echo "  schema           Regenerate schema.json and bridge-metadata.json"
	@echo "  build_sdks       Regenerate sdk/{nodejs,python,go,dotnet,java}"
	@echo "  build_nodejs     Regenerate sdk/nodejs only"
	@echo "  install          go install the runtime provider into GOPATH/bin"
	@echo "  fmt              gofumpt + goimports in provider/"
	@echo "  lint             golangci-lint run in provider/"
	@echo "  lint_fix         golangci-lint run --fix in provider/"
	@echo "  tidy             go mod tidy in provider/"
	@echo "  clean            Remove bin/ and any sdk-tmp/"
	@echo ""
	@echo "Current VERSION: $(VERSION)"

all: build

build: schema provider build_sdks

# --- Go builds ---------------------------------------------------------

$(BIN):
	@mkdir -p $(BIN)

tfgen: | $(BIN)
	cd $(PROVIDER_PATH) && go build -o $(BIN)/$(TFGEN) -ldflags "$(LDFLAGS)" ./cmd/$(TFGEN)

provider: schema | $(BIN)
	cd $(PROVIDER_PATH) && go build -o $(BIN)/$(PROVIDER) -ldflags "$(LDFLAGS)" ./cmd/$(PROVIDER)

# Run tfgen to (re)generate the Pulumi schema + bridge metadata. Both get
# committed to git so consumers don't need the toolchain to inspect them.
schema: tfgen
	$(BIN)/$(TFGEN) schema --out $(PROVIDER_PATH)/cmd/$(PROVIDER)

# --- SDK generation ----------------------------------------------------

# SDK generation uses the Pulumi CLI's built-in codegen, not the bridge's
# tfgen. tfgen only emits schema.json; `pulumi package gen-sdk` handles
# the actual TypeScript / Python / Go / .NET / Java rendering.

build_nodejs: schema
	pulumi package gen-sdk $(SCHEMA_PATH) --language nodejs --out sdk-tmp
	rm -rf sdk/nodejs
	mv sdk-tmp/nodejs sdk/nodejs
	rm -rf sdk-tmp

build_python: schema
	pulumi package gen-sdk $(SCHEMA_PATH) --language python --out sdk-tmp
	rm -rf sdk/python
	mv sdk-tmp/python sdk/python
	rm -rf sdk-tmp

build_go: schema
	pulumi package gen-sdk $(SCHEMA_PATH) --language go --out sdk-tmp
	rm -rf sdk/go
	mv sdk-tmp/go sdk/go
	rm -rf sdk-tmp

build_dotnet: schema
	pulumi package gen-sdk $(SCHEMA_PATH) --language dotnet --out sdk-tmp
	rm -rf sdk/dotnet
	mv sdk-tmp/dotnet sdk/dotnet
	rm -rf sdk-tmp

build_java: schema
	pulumi package gen-sdk $(SCHEMA_PATH) --language java --out sdk-tmp
	rm -rf sdk/java
	mv sdk-tmp/java sdk/java
	rm -rf sdk-tmp

build_sdks: build_nodejs build_python build_go build_dotnet build_java

# --- Installation ------------------------------------------------------

install: install_provider

install_provider: provider
	install -m 0755 $(BIN)/$(PROVIDER) $(shell go env GOPATH)/bin/

install_nodejs_sdk:
	cd sdk/nodejs && npm install && npm run build && cd bin && npm pack

# --- Dev ergonomics ----------------------------------------------------

fmt:
	cd $(PROVIDER_PATH) && gofumpt -w . && goimports -w -local $(PROJECT) .

lint:
	cd $(PROVIDER_PATH) && golangci-lint run ./...

lint_fix:
	cd $(PROVIDER_PATH) && golangci-lint run --fix ./...

tidy:
	cd $(PROVIDER_PATH) && go mod tidy

clean:
	rm -rf $(BIN) sdk-tmp
