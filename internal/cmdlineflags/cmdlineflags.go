package cmdlineflags

import (
	"flag"
)

// FlagsRead - Wrapper for command line flags read
type FlagsRead struct {
	// Full XML path
	BookPath string
	// Whether to split images into their own xhtml file
	SplitImages bool
}

func ParseCommandLineFlags() FlagsRead {
	bookPath := flag.String("bookPath", "", "Full path of .xhtml file")
	splitImages := flag.Bool("splitImages", true, "Whether to split images into their own xhtml file")
	flag.Parse()

	if bookPath == nil || len(*bookPath) == 0 {
		panic("Please supply a `bookPath` value")
	}

	if splitImages == nil {
		panic("Please supply a `splitImages` value")
	}

	return FlagsRead{
		BookPath:    *bookPath,
		SplitImages: *splitImages,
	}
}
