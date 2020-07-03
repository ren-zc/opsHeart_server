.PHONY: build

BUILDTIME=$(shell date +%FT%T%z)
GOVERSION=$(shell go version)

all: build

build:
	go build -ldflags "-X 'main.buildTime=${BUILDTIME}' -X 'main.goVersion=${GOVERSION}'" -o ./super_server .
