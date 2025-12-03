VERSION := $(shell git describe --tags --abbrev=0)
COMMIT := $(shell git rev-parse --short HEAD)

build:
  go build -ldflags "\
    -w -s \
    -X main.Version=$(VERSION) \
    -X main.Commit=$(COMMIT)" \
    -o ./bin/monoctl ./main.go
