package main

import (
	"log"
	"os"
	"runtime/pprof"
	"runtime/trace"
)

func cpuProfiling(fn string) {

	log.Println("Starting profiling")
	f, err := os.Create(fn)
	if err != nil {
		log.Fatal(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
}

func cpuTracing(fn string) {

	log.Println("Starting Tracing")
	f, err := os.Create(fn)
	if err != nil {
		log.Fatal(err)
	}
	trace.Start(f)
	defer trace.Stop()
}
