VERSION := $(shell git describe --tags --abbrev=0)
COMMIT := $(shell git rev-parse --short HEAD)

rm:
	rm -rf ./bin

build: rm
	go build -ldflags " \
		-w -s \
		-X 'github.com/salmanwaheed/monoctl.BuildInfo.Version=$(VERSION)' \
		-X 'github.com/salmanwaheed/monoctl.BuildInfo.Commit=$(COMMIT)'" \
		-o ./bin/monoctl ./cmd/monoctl/main.go
