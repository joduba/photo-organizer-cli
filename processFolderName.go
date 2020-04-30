package main

import (
	"regexp"
	"time"

	"github.com/gosimple/slug"
)

// ProcessFolder gets an string and will try to extract a date and a suffix from it
func ProcessFolder(folder string) (string, time.Time) {
	/*
	   S: 2020-04-30 ... ["2020-04-30" "2020" "04" "30" "" ""]
	   S: 2020-04-30 Park Güell ... ["2020-04-30 Park Güell" "2020" "04" "30" "" "Park Güell"]
	   S: 2020-04-30-Park Güell ... ["2020-04-30-Park Güell" "2020" "04" "30" "" "Park Güell"]
	   S: 2020-04-30-Park Güell ... ["2020-04-30-Park Güell" "2020" "04" "30" "" "Park Güell"]
	   S: 2020 04 30 Park Güell ... ["2020 04 30 Park Güell" "2020" "04" "30" "" "Park Güell"]
	   S: 2020/04/30 Park Güell ... ["2020/04/30 Park Güell" "2020" "04" "30" "" "Park Güell"]
	   S: 2020-04-22-30 Park Güell ... ["2020-04-22-30 Park Güell" "2020" "04" "22" "30" "Park Güell"]
	   S: 2020-04 Park Güell ... ["2020-04 Park Güell" "2020" "04" "" "" "Park Güell"]
	   S: Park Güell ... ["Park Güell" "" "" "" "" "Park Güell"]
	*/

	re := regexp.MustCompile(`([0-9]{4})?[\-\ \/]?([0-9]{2})?[\-\ \/]?([0-9]{2})?[\-\ \/]?([0-9]{2})?[\-\ ]*(.*)`)
	p := re.FindStringSubmatch(folder)

	s := slug.Make(p[5])

	if p[1] == "" {
		return s, time.Time{}
	}

	f := "2006"
	d := p[1]
	if p[2] != "" {
		f += "-" + "01"
		d += "-" + p[2]
	}
	if p[3] != "" {
		f += "-" + "02"
		d += "-" + p[3]
	}

	t, _ := time.Parse(f, d)

	return s, t
}
