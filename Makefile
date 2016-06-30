PKG?=$(shell go list ./... | grep -v /vendor/)
VETARGS?=-asmdecl -atomic -bool -buildtags -copylocks -methods -nilfunc -printf -rangeloops -shift -structtags -unsafeptr

default: build

build:
	go install -v github.com/brettbuddin/ponyexpress/cmd/ponyexpress

vet:
	go vet $(PKG)

cover:
	go test $(PKG) -coverprofile=coverage.out
	go tool cover -html=coverage.out
	@rm coverage.out

test:
	go test $(PKG) $(TESTARGS)

.PHONY: default test build cover vet benchmark
