package main

import (
	"flag"
	"log"
	"os"
	"time"
)

func main() {
	var suffix string
	var timeOffset int

	log.Println("photo: Starting")
	flag.IntVar(&timeOffset, "offset", 0, "Number of hours to be added (or removed) to the current time")
	flag.StringVar(&suffix, "suffix", "", "Text to be added to enrich the name of the files")
	flag.Parse()

	if flag.Arg(0) == "" {
		log.Fatal("no folder provided")
		os.Exit(1)
	}

	log.Printf("-- Parameters: suffix [%s], offset [%d] folder [%v]\n", suffix, timeOffset, flag.Arg(0))

	startTime := time.Now()

	c, err := renameAndChangeTime(flag.Arg(0), suffix, timeOffset)
	endTime := time.Now()
	etime := endTime.Sub(startTime)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	if c == 0 {
		log.Fatal("no photos to rename in %s", etime)
	}
	log.Printf("Successfully renamed %d photos in %s\n", c, etime)

}