# Code snippets & references
Those are a bunch of code snipets from other repos and example that I may want to use some day

## From "github.com/rwcarlsen/goexif/exif"

```
{
	exif.Walk(Walker{})

}

	type Walker struct{}

func (_ Walker) Walk(name exif.FieldName, tag *tiff.Tag) error {
	data, _ := tag.MarshalJSON()
	fmt.Printf("    %v: %v\n", name, string(data))
	return nil
}
```

## From: "github.com/tajtiattila/metadata"

```
m, err := metadata.ParseAt(f)
	if err != nil {
		fmt.Printf("----SKIP: file: %s\n", path)
		//log.Printf("ERROR: parsing file %s: %s", path,err)
		return nil
	}

	fmt.Printf("META - FILE: %s:\n", path)
	for k, v := range m.Attr {
		fmt.Printf("META  %s: %q\n", k, v)
	}
```
