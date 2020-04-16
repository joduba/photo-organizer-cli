# Go related variables.
GOBIN := $(GOPATH)/bin
GOFILE := $(shell basename "$(PWD)")
GOPPROF := $(shell basename "$(PWD)").prof

build:
	go build -o $(GOFILE) *.go

pprof:
	make build
	./$(GOFILE) --pprof=$(GOPPROF) workdir
	go tool pprof $(GOFILE) $(GOPPROF) 