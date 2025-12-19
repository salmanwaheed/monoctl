VERSION := $(shell git describe --tags --abbrev=0)
GIT_COMMIT := $(shell git rev-parse --short HEAD)
GO_VERSION := $(shell go version | awk '{print $$3}')
OS := $(shell uname --kernel-name | tr 'A-Z' 'a-z')
ARCH := $(shell uname --machine)

rm:
	rm -rf ./bin

build: rm
	go build -ldflags " \
		-w -s \
		-X 'github.com/salmanwaheed/monoctl.BuildInfo.Version=$(VERSION)' \
		-X 'github.com/salmanwaheed/monoctl.BuildInfo.GitCommit=$(GIT_COMMIT)' \
		-X 'github.com/salmanwaheed/monoctl.BuildInfo.GoVersion=$(GO_VERSION)' \
		-X 'github.com/salmanwaheed/monoctl.BuildInfo.OS=$(OS)' \
		-X 'github.com/salmanwaheed/monoctl.BuildInfo.Arch=$(ARCH)'" \
		-o ./bin/monoctl-$(OS)-$(ARCH) ./cmd/monoctl/main.go
