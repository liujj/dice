APPNAME = dice
TAG = $(shell git describe --tags `git rev-list --tags --max-count=1`)
LDFLAGS += -X "dice/pkg/version.BuildTime=$(shell date "+%Y-%m-%d %T %Z")"
LDFLAGS += -X "dice/pkg/version.GitCommit=$(shell git rev-parse HEAD)"
LDFLAGS += -X "dice/pkg/version.LatestTag=$(TAG)"

all: linux out package clean

.PHONY: local
local:
	CGO_ENABLED=1 go build -ldflags '$(LDFLAGS)' -o ./bin/$(APPNAME) ./cmd/main.go

.PHONY: linux
linux:
	CGO_ENABLED=1 GOOS=linux go build -ldflags '$(LDFLAGS)' -o ./bin/$(APPNAME) ./cmd/main.go

.PHONY: rpi
linux:
	CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -ldflags '$(LDFLAGS)' -o ./bin/$(APPNAME) ./cmd/main.go

.PHONY: out
out:
	@if [ -e out ] ; then rm -rf out ; fi
	@mkdir out
	@cp ./bin/$(APPNAME) ./out
	@cp ./config/config.json ./out
	@cp ./db.sqlite3 ./out

.PHONY: package
package:
	@if [ -e package ] ; then rm -rf package ; fi
	@mkdir package
	@cd ./out && tar -zcf ../package/$(APPNAME).$(TAG).tar.gz ./*

.PHONY: package-all
package-all:
	@if [ -e $(APPNAME).tar.gz ] ; then rm $(APPNAME).tar.gz ; fi
	@cd ./out && tar -zcf ../$(APPNAME).tar.gz ./*

.PHONY: clean
clean:
	@rm -rf ./bin
	@rm -rf ./out
