VERSION := $(shell git describe --tags --abbrev=0)
COMMIT := $(shell git rev-parse --short HEAD)

build:
	go build -ldflags " \
		-w -s \
		-X 'github.com/salmanwaheed/monoctl/cmd.VersionInfo.Version=$(VERSION)' \
		-X 'github.com/salmanwaheed/monoctl/cmd.VersionInfo.Commit=$(COMMIT)'" \
		-o ./bin/monoctl ./main.go
