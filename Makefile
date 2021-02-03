PWD := $(shell pwd)
GOBIN ?= $(PWD)/bin
GO_OPTS ?= GOBIN=$(GOBIN)

GLOBAL_GOPATH := $(GOPATH)
LOCAL_GOPATH  := $(HOME)/go

ifdef GLOBAL_GOPATH
  GOPATH := $(GLOBAL_GOPATH)
else
  GOPATH := $(LOCAL_GOPATH)
endif

install: # Install `evolve`.
	$(GO_OPTS) go install $(PWD)/evolve/
