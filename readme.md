# PhotoOrganizer

This is an simple go application that will read the pictures on a folder and rename them based on the date the picture was taken from EXIF data. While doing the rename, the application will also set the access and modified time of the file to the same date.

Optionally, you can specify a suffix to be added to the pictures

`--suffix "my suffix"`

And also you can specify a time offset in hours. This is useful for cameras that we have forgot to adjust to summer time or when traveling we did not change to the local time and now we want the pictures to reflect the time on that location not on our country of origin.

`--offset 1`

This is a work in progress and some of the features to be added are:

- move the organized folder to an specified folder
- classify the pictures inside a global folder by topic (user defined)
- automatically detect the year and organize move the folder to that year inside the global folder/topic
- load a configuration file where to store common configuration
- Tensorflow image tagging
- update XMP or EXIF metadata with tags and new date in case of time change.

## Building the application

`go build -o photo-organizer main.go rename-and-chtime.go `

## Examples

Get help

`photo-organizer --help`

Rename all image files on the folder `workdir` and add the suffix "weekend-trip" to it.

`photo-organizer --suffix "weekend-trip" workdir`

Rename all image files on the folder `workdir` and change the date by minus 1h to adjust to the real time we toke the picture.

`photo-organizer --offset -1h  workdir`

