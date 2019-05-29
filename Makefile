GOPATH := $(shell dirname $(shell dirname $(shell dirname $(shell dirname $(shell realpath .)))))
NS := github.com/projecteru2/yavirt
GO := GOPATH=$(GOPATH) go
BUILD := $(GO) build -race
TEST := $(GO) test -count=1 -race -cover

LDFLAGS += -X "$(NS)/ver.Git=$(shell git rev-parse HEAD)"
LDFLAGS += -X "$(NS)/ver.Compile=$(shell $(GO) version)"
LDFLAGS += -X "$(NS)/ver.Date=$(shell date +'%F %T %z')"

PKGS := $$($(GO) list ./... | grep -v vendor/)

.PHONY: all test build

default: build

build: lint
	$(BUILD) -ldflags '$(LDFLAGS)' -o bin/yavirtd yavirtd.go

lint: format

format:
	$(GO) vet $(PKGS)
	$(GO) fmt $(PKGS)

deps:
	GO111MODULE=on GOPATH=$(GOPATH) go mod download
	GO111MODULE=on GOPATH=$(GOPATH) go mod vendor

clean:
	rm -fr bin/*

test:
ifdef RUN
	$(TEST) -v -run='${RUN}' $(PKGS)
else
	$(TEST) $(PKGS)
endif

initdev:
	mysql -e 'DROP DATABASE test; CREATE DATABASE test'
	$(GO) run yavirtd.go --init
	mysql test < schema/init.sql
	mysql test -e "INSERT host_tab (hostname, subnet) SELECT '`hostname`', 12625921"
	mysql test -e "INSERT addr_tab (low_value, prefix, gateway, state, host_subnet) SELECT 3232235777, 24, '192.168.122.254', 'free', 12625921"

rundev:
	sudo GOPATH=$(GOPATH) $(GO) run yavirtd.go
