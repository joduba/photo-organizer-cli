//This is an simple go application that will read the pictures on a folder and rename them based on the date the picture was taken from EXIF data.
//While doing the rename, the application will also set the access and modified time of the file to the same date.
package main

import (
	"flag"
	"log"
	"os"
	"time"
)

func main() {
	var suffix, cpuprofile, cputrace string
	var timeOffset int

	log.Println("photo: Starting")
	flag.IntVar(&timeOffset, "offset", 0, "Number of hours to be added (or removed) to the current time")
	flag.StringVar(&suffix, "suffix", "", "Text to be added to enrich the name of the files")
	flag.Parse()

	if flag.Arg(0) == "" {
		log.Fatal("no folder provided")
		os.Exit(1)
	}

	log.Printf("-- Parameters: picture suffix [%s], time offset [%d Hours] folder [%v]\n", suffix, timeOffset, flag.Arg(0))

	startTime := time.Now()

	c, err := renameAndChangeTime(flag.Arg(0), suffix, timeOffset)
	endTime := time.Now()
	etime := endTime.Sub(startTime)

	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	if c == 0 {
		log.Fatalf("no photos to rename in %s", etime)
	}
	log.Printf("Successfully renamed %d photos in %s\n", c, etime)

}
