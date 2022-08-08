package cmdlineflags

import (
	"flag"
	"strings"
)

// FlagsRead - Wrapper for command line flags read
type FlagsRead struct {
	// What command was issued
	Command string
	// Full XML path
	BookPath string
	// Whether to split images into their own xhtml file
	SplitImages bool
}

func ParseCommandLineFlags() FlagsRead {
	command := flag.String("command", "", "Valid commands are: "+strings.Join(ValidCommmands, ", "))
	bookPath := flag.String("bookPath", "", "Full path of .xhtml file")
	splitImages := flag.Bool("splitImages", true, "Whether to split images into their own xhtml file")
	flag.Parse()

	if command == nil || len(*command) == 0 {
		panic("Please supply command")
	}

	if bookPath == nil || len(*bookPath) == 0 {
		panic("Please supply a `bookPath` value")
	}

	if splitImages == nil {
		panic("Please supply a `splitImages` value")
	}

	return FlagsRead{
		Command:     *command,
		BookPath:    *bookPath,
		SplitImages: *splitImages,
	}
}
