# PhotoOrganizer

This is an simple go application that will read the pictures on a folder and rename them based on the date the picture was taken from EXIF data. While doing the rename, the application will also set the access and modified time of the file to the same date.

## Options 

`-suffix "my suffix"`

Optionally, you can specify a suffix to be added to the pictures

`-offset 1`

And also you can specify a time offset in hours. This is useful for cameras that we have forgot to adjust to summer time or when traveling we did not change to the local time and now we want the pictures to reflect the time on that location not on our country of origin.

`-classify`

There is also the possibility to classify a group of picture from different dates into a folder structure splitting by years and days.

`-auto`

And if you want to integrate with MacOs automator, it will be useful to get all the information from the folder name. With the option `-auto` will take the folder name and look for a date patern that will use as default date for the pictures that has no exif (like the ones saved from Facebook), and the rest of the folder name will be used as suffix. If you are interested in this, check the folder automator, there are a couple of examples you can use in your mac.

## Building the application

`go build -o photoOrganizer main.go rename-and-chtime.go`

## Examples

Get help

`photoOrganizer -help`

Rename all image files on the folder `workdir` and add the suffix "weekend-trip" to it.

`photo-organizer -suffix "weekend-trip" workdir`

Rename all image files on the folder `workdir` and change the date by minus 1h to adjust to the real time we toke the picture.

`photo-organizer -offset -1h workdir`

Create a new folder structure under the folder `basedir`, organized by `year/year-month-day-suffix` and rename and move all image files present at the `workdir` folder. If we ommit the basedir, it will create a folder called `out`.

`photo-organizer -classify -basedir "myPhotos" -suffix "red-moon" workdir`

Guess the suffix from the foldername and use if needed also the default date. It will rename the folder with the format YYYY-MM-DAY-Suffix. The Suffix will be converted into an Slug to avoid issues. Works with and without `classify`

`photo-organizer -auto -classify workdir`
