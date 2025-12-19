VERSION := $(shell git describe --tags --abbrev=0)
GIT_COMMIT := $(shell git rev-parse --short HEAD)
GO_VERSION := $(shell go version | awk '{print $$3}')
KERNEL_NAME := $(shell uname --kernel-name | tr 'A-Z' 'a-z')
ARCH := $(shell uname --machine)

rm:
	rm -rf ./bin

build: rm
	go build -ldflags " \
		-w -s \
		-X 'github.com/salmanwaheed/monoctl.version=$(VERSION)' \
		-X 'github.com/salmanwaheed/monoctl.gitCommit=$(GIT_COMMIT)' \
		-X 'github.com/salmanwaheed/monoctl.goVersion=$(GO_VERSION)' \
		-X 'github.com/salmanwaheed/monoctl.kernelName=$(KERNEL_NAME)' \
		-X 'github.com/salmanwaheed/monoctl.arch=$(ARCH)'" \
		-o ./bin/monoctl-$(KERNEL_NAME)-$(ARCH) ./cmd/monoctl/main.go
