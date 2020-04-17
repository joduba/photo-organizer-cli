# Go related variables.
GOBIN := $(GOPATH)/bin
GOFILE := $(shell basename "$(PWD)")
GOPPROF := $(shell basename "$(PWD)").prof
GOTRACE := $(shell basename "$(PWD)").trace

build:
	go build -o $(GOFILE) *.go
