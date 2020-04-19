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
	var suffix, basedir string
	var timeOffset int
	var classify, auto bool

	log.Println("photo: Starting")
	flag.IntVar(&timeOffset, "offset", 0, "Number of hours to be added (or removed) to the current time")
	flag.StringVar(&suffix, "suffix", "", "Text to be added to enrich the name of the files")
	flag.BoolVar(&classify, "classify", false, "If set to true, it will organize the pictures in folders by year/day")
	flag.BoolVar(&auto, "auto", false, "If set to true, will use the folder name as Suffix")
	flag.StringVar(&basedir, "basedir", "", "Base folder where move the files and folders if classify flag is set. Default is inside the work folder")
	flag.Parse()

	startTime := time.Now()
	log.Printf("-- Parameters: picture suffix [%s], time offset [%d Hours] folder [%v] baseDir [%s]\n", suffix, timeOffset, flag.Arg(0), basedir)

	folder := flag.Arg(0)
	if folder == "" {
		log.Fatal("no folder provided")
		os.Exit(1)
	}
	fileInfo, err := os.Lstat(folder)
	if err != nil {
		log.Fatal("can't STAT the folder")
		os.Exit(2)
	}
	if !fileInfo.IsDir() {
		log.Fatal("The provided argument is not a folder")
		os.Exit(3)
	}

	if classify == false && basedir != "" {
		log.Printf("Igonring -basedir '%s', it's only used when working with -classify\n", basedir)
	}

	if basedir == "" && classify == true {
		basedir = folder
	}

	if auto == true && suffix != "" {
		suffix = folder
	}

	c, err := doFolderRename(folder, suffix, timeOffset, classify, basedir)
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
