LDFLAGS += -X "coinsBook/version.BuildTime=$(shell date "+%Y-%m-%d %T %Z")"
LDFLAGS += -X "coinsBook/version.GitCommit=$(shell git rev-parse HEAD)"
OS := $(shell uname -s).$(shell uname -m)
GOVET = go tool vet -composites=false -methods=false -structtags=false
GOFMT ?= gofmt "-s"
GOFILES := $(shell find . -name "*.go" -type f -not -path "./vendor/*")
PACK_FILE = blockchain_coinsbook
NOW = $(shell date -u '+%Y%m%d%I%M%S')
.PHONY: build pack clean out


all: build out pack clean

build:
	@go build -ldflags '$(LDFLAGS)' -o ./bin/tm_app ./main.go
	@go build -ldflags '$(LDFLAGS)' -o ./bin/tm_launcher ./launcher/main.go

out:
	@if [ -e out ] ; then rm -rf out; fi
	@mkdir out
	@cp ./bin/tm_app ./out
	@cp ./bin/tm_launcher ./out
	@cp ./config/config.json ./out
	@cp -r ./view ./out

pack:
	@cd out;zip -rq9 $(PACK_FILE).zip .;cd ..
	@mv ./out/$(PACK_FILE).zip .

clean:
	@rm -rf ./bin
	@rm -rf ./out
	@go clean

cleanp:
	clean
	@rm -f $(PACK_FILE).tar.gz

fmt:
	@$(GOFMT) -w $(GOFILES)