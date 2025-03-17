
PACKAGES=$(shell go list ./...)
BUILDDIR?=$(CURDIR)/build
OUTPUT?=$(BUILDDIR)/cometbft
VALIDATORS_COUNT?=4

HTTPS_GIT := https://github.com/cometbft/cometbft.git
CGO_ENABLED ?= 0

# Process Docker environment variable TARGETPLATFORM
# in order to build binary with correspondent ARCH
# by default will always build for linux/amd64
TARGETPLATFORM ?=
GOOS ?= linux
GOARCH ?= amd64
GOARM ?=

ifeq (linux/arm,$(findstring linux/arm,$(TARGETPLATFORM)))
	GOOS=linux
	GOARCH=arm
	GOARM=7
endif

ifeq (linux/arm/v6,$(findstring linux/arm/v6,$(TARGETPLATFORM)))
	GOOS=linux
	GOARCH=arm
	GOARM=6
endif

ifeq (linux/arm64,$(findstring linux/arm64,$(TARGETPLATFORM)))
	GOOS=linux
	GOARCH=arm64
	GOARM=7
endif

ifeq (linux/386,$(findstring linux/386,$(TARGETPLATFORM)))
	GOOS=linux
	GOARCH=386
endif

ifeq (linux/amd64,$(findstring linux/amd64,$(TARGETPLATFORM)))
	GOOS=linux
	GOARCH=amd64
endif

ifeq (linux/mips,$(findstring linux/mips,$(TARGETPLATFORM)))
	GOOS=linux
	GOARCH=mips
endif

ifeq (linux/mipsle,$(findstring linux/mipsle,$(TARGETPLATFORM)))
	GOOS=linux
	GOARCH=mipsle
endif

ifeq (linux/mips64,$(findstring linux/mips64,$(TARGETPLATFORM)))
	GOOS=linux
	GOARCH=mips64
endif

ifeq (linux/mips64le,$(findstring linux/mips64le,$(TARGETPLATFORM)))
	GOOS=linux
	GOARCH=mips64le
endif

ifeq (linux/riscv64,$(findstring linux/riscv64,$(TARGETPLATFORM)))
	GOOS=linux
	GOARCH=riscv64
endif

install-cometbft:
	go install github.com/cometbft/cometbft/cmd/cometbft@v1.0

build: clean
	go mod tidy
	mkdir -p build
	CGO_ENABLED=$(CGO_ENABLED) go build $(BUILD_FLAGS) -tags '$(BUILD_TAGS)' -o $(OUTPUT)
	cp ./localnode/config-template.toml $(BUILDDIR)/config-template.toml

init:
	cometbft init --home /tmp/cometbft-home

start:
	./test -cmt-home /tmp/cometbft-home

build-docker-localnode:
	docker buildx build --no-cache --platform linux/amd64 --tag cometbft/localnode localnode

start-localnet: stop-localnet  build-docker-localnode
	@if ! [ -f build/node0/config/genesis.json ]; then docker run --rm -v $(CURDIR)/build:/cometbft cometbft/cometbft:v1.0.0 testnet --config ./config-template.toml --v $(VALIDATORS_COUNT) --o . --populate-persistent-peers --starting-ip-address 192.167.10.2 ; fi
	docker compose up -d

stop-localnet:
	docker compose down

clean:
	rm -rf build

build-linux:
	GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) $(MAKE) build

apply-latency:
	./latency.sh $(VALIDATORS_COUNT)