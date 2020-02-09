APPNAME = dice
TAG = $(shell git describe --tags `git rev-list --tags --max-count=1`)
LDFLAGS += -X "pkg/dice/version.BuildTime=$(shell date "+%Y-%m-%d %T %Z")"
LDFLAGS += -X "pkg/dice/version.GitCommit=$(shell git rev-parse HEAD)"
LDFLAGS += -X "pkg/dice/version.LatestTag=$(TAG)"

all: depend linux out package clean

.PHONY: depend
depend:
	go get github.com/rakyll/statik

.PHONY: statik
statik:
	statik -src=./view -dest=./pkg

.PHONY: local
local: statik
	CGO_ENABLED=1 go build -ldflags '$(LDFLAGS)' -o ./bin/$(APPNAME) ./cmd/main.go

.PHONY: linux
linux: statik
	CGO_ENABLED=1 GOOS=linux go build -ldflags '$(LDFLAGS)' -o ./bin/$(APPNAME) ./cmd/main.go

.PHONY: out
out:
	@if [ -e out ] ; then rm -rf out ; fi
	@mkdir out
	@cp ./bin/$(APPNAME) ./out
	@cp ./config/config.json ./out
	@cp -r ./view ./out

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
