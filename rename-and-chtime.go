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

func renameAndChangeTime(f, s string, o int) (int, error) {

	fmt.Printf("-- Renaming images in folder '%s' with suffix [%s] and offset [%d]...\n", f, s, o)
	timeOffset = time.Duration(o) * time.Hour
	suffix = s

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
	for c := 1; true; c++ {
		new := fmt.Sprintf(t, c)
		if name == new {
			return path + "/" + new, nil
		}
		if _, err := os.Stat(path + "/" + new); err != nil {
			return path + "/" + new, nil
		}
	}
	return "", fmt.Errorf("x Could not find available filename for: %s", name)
}
