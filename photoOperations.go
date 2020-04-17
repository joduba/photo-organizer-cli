package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

var extensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".nef":  true,
	".heic": true,
}

var changed = 0
var suffix = ""
var timeOffset time.Duration
var classify = false
var dstFolders map[string]bool
var basedir = ""

func doPhotoOperations(f, s string, o int, c bool, b string) (int, error) {

	fmt.Printf("-- Renaming images in folder '%s' with suffix [%s] and offset [%d]...\n", f, s, o)
	timeOffset = time.Duration(o) * time.Hour
	suffix = s
	basedir = b

	classify = c
	if classify {
		dstFolders = make(map[string]bool)
	}

	err := filepath.Walk(f, doRenameAndChangeTime)
	if err != nil {
		return changed, err
	}

	fmt.Println()
	return changed, nil

}

func doRenameAndChangeTime(path string, fileInfo os.FileInfo, err error) error {

	if err != nil {
		return fmt.Errorf("an error ocurred while accessing a path %q:\n %v", path, err)
	}

	if fileInfo == nil {
		return fmt.Errorf("while trying to read the file %q:\n %v", path, err)
	}

	// is a dir, nothing to do
	if fileInfo.IsDir() {
		return nil
	}

	//not a foto
	if extension := strings.ToLower(filepath.Ext(path)); extensions[extension] != true {
		//log.Printf("*** extension %s is not a photo\n", extension)
		return nil
	}

	f, _ := os.Open(path)
	defer f.Close()
	if err != nil {
		return fmt.Errorf("opening file: %s\n %v", path, err)
	}

	pDateTime := fileInfo.ModTime()
	if exif, err := exif.Decode(f); err != nil {
		pDateTime, _ = exif.DateTime()
	}

	if timeOffset != 0 {
		pDateTime = pDateTime.Add(timeOffset)
	}

	newName, err := findName(filepath.Dir(path), fileInfo.Name(), pDateTime, suffix)
	if err != nil {
		return fmt.Errorf("can find a new name for %s\n %v", path, err)
	}

	// No need to rename if the new name == oldname
	if path != newName {
		if err = os.Rename(path, newName); err != nil {
			return fmt.Errorf("can't rename the file %s\n %v", path, err)
		}
	}

	if err := os.Chtimes(newName, pDateTime, pDateTime); err != nil {
		return fmt.Errorf("change time issue for file: %s\n %v", path, err)
	}
	changed++
	log.Printf("✓ Renamed from %s to: %s and CHTimes to %v\n", path, newName, pDateTime)

	return nil
}

func findName(path, name string, pDateTime time.Time, suffix string) (string, error) {
	ftime := pDateTime.Format("20060102-150405")
	t := ftime + "-%03d"
	if suffix != "" {
		t += "-" + suffix
	}
	t += filepath.Ext(name)
	dstPath := path
	var err error
	if classify == true {
		if dstPath, err = destinationPath(pDateTime); err != nil {
			return "", err
		}
	}
	for c := 1; true; c++ {
		new := dstPath + "/" + fmt.Sprintf(t, c)
		if path+"/"+name == new {
			return new, nil
		}
		if _, err := os.Stat(new); err != nil {
			return new, nil
		}
	}
	return "", fmt.Errorf("x Could not find available filename for: %s", name)
}

func destinationPath(d time.Time) (string, error) {

	f := basedir + "/" + d.Format("2006/2006-01-02")
	if suffix != "" {
		f += "-" + suffix
	}
	if _, ok := dstFolders[f]; ok {
		return f, nil
	}

	if err := os.MkdirAll(f, 0777); err != nil {
		return "", fmt.Errorf("x Could not create folder: %s", f)
	}
	dstFolders[f] = true

	return f, nil
}
