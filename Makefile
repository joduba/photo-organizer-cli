# Go related variables.
GOBIN := $(GOPATH)/bin
GOFILE := $(shell basename "$(PWD)")
GOPPROF := $(shell basename "$(PWD)").prof
GOTRACE := $(shell basename "$(PWD)").trace

build:
	go build -o $(GOFILE) *.go

pprof:
	make build
	./$(GOFILE) --pprof=$(GOPPROF) workdir
	go tool pprof $(GOFILE) $(GOPPROF) 

trace:
	make build
	./$(GOFILE) --trace=$(GOTRACE) workdir
	go tool trace $(GOTRACE)