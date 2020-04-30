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
var suffix, basedir string
var timeOffset time.Duration
var classify = false
var dstFolders map[string]bool
var defaultDate time.Time

func doFolderRename(f, s string, o int, c bool, b string, a bool) (int, error) {

	fmt.Printf("-- Renaming images in folder '%s' with suffix [%s] and offset [%d]...\n", f, s, o)
	timeOffset = time.Duration(o) * time.Hour
	suffix = s
	if b != "" {
		basedir = b + "/"
	}

	classify = c
	dstFolders = make(map[string]bool)

	fi, err := os.Stat(f)
	if err != nil {
		return changed, fmt.Errorf("an error ocurred while accessing the main folder path %q:\n %v", f, err)
	}

	if fi.IsDir() == false {
		return changed, fmt.Errorf("folder provide is not a folder %q:\n %v", f, err)
	}

	if a == true {
		suffix, defaultDate = ProcessFolder(fi.Name())
		log.Printf(" Auto is on, detected suffix: [%s] and default date: [%v]\n", suffix, defaultDate)
	}

	err = filepath.Walk(f, doPhotoRename)
	if err != nil {
		return changed, err
	}

	if a == false {
		return changed, nil
	}

	err = folderRename(f, fi, defaultDate, suffix)
	if err != nil {
		return changed, err
	}

	return changed, nil

}

func doPhotoRename(path string, fileInfo os.FileInfo, err error) error {

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
	if exif, err := exif.Decode(f); err == nil && exif != nil {
		pDateTime, _ = exif.DateTime()
	} else {
		log.Printf("-- Weird, the picture %s does not have a valid exif, using fileInfo.ModTime instead \n", path)
	}

	if timeOffset != 0 {
		pDateTime = pDateTime.Add(timeOffset)
	}

	newName, err := findFileName(filepath.Dir(path), fileInfo.Name(), pDateTime, suffix)
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
	log.Printf("✓ renamed from %s to: %s and CHTimes to %v\n", path, newName, pDateTime)

	return nil
}

func folderRename(path string, fi os.FileInfo, pDateTime time.Time, suffix string) error {

	newName, err := findFolderName(path, fi.Name(), pDateTime, suffix)
	if err != nil {
		return fmt.Errorf("can find a new name for main folder %s\n %v", path, err)
	}
	// No need to rename if the new name == oldname
	if path != newName {
		if err = os.Rename(path, newName); err != nil {
			return fmt.Errorf("can't rename the folder %s\n %v", path, err)
		}
		log.Printf("✓ renamed main folder %v -> %v\n", path, newName)
	}
	return nil
}

func findFolderName(path, name string, pDateTime time.Time, suffix string) (string, error) {

	parentFolder := filepath.Dir(path)
	t := pDateTime.Format("20060102")
	if suffix != "" {
		t += "-" + suffix
	}
	t += filepath.Ext(name)

	log.Printf(" +++++ parentFolder: [%s], name = [%s], new [%s], path[%s]\n", parentFolder, name, t, path)

	if name == t {
		return t, nil
	}

	t = parentFolder + "/" + t
	if _, err := os.Stat(t); err != nil {
		return t, nil
	}

	// the name already exist, and we need to do something.
	t = parentFolder + "/" + pDateTime.Format("20060102")
	if suffix != "" {
		t += "-" + suffix
	}
	t += "-%02d" + filepath.Ext(name)

	for c := 1; true; c++ {
		new := fmt.Sprintf(t, c)
		if path == new {
			return new, nil
		}
		if _, err := os.Stat(new); err != nil {
			return new, nil
		}
	}
	return "", fmt.Errorf("x Could not find available filename for: %s", name)
}

func findFileName(path, name string, pDateTime time.Time, suffix string) (string, error) {
	t := pDateTime.Format("20060102-150405") + "-%03d"
	if suffix != "" {
		t += "-" + suffix
	}
	t += filepath.Ext(name)
	dstPath := ""
	var err error
	if dstPath, err = destinationPath(pDateTime, path); err != nil {
		return "", err
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

func destinationPath(d time.Time, path string) (string, error) {

	if classify != true {
		return path, nil
	}
	f := basedir + d.Format("2006/2006-01-02")
	if suffix != "" {
		f += "-" + suffix
	}

	if _, ok := dstFolders[f]; ok {
		return f, nil
	}

	if err := os.MkdirAll(f, 0775); err != nil {
		return "", fmt.Errorf("x Could not create folder: %s", f)
	}
	dstFolders[f] = true

	return f, nil
}
