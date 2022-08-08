package main

import (
	"strings"

	"github.com/toddnguyen47/splith1headers/internal/cmdlineflags"
	"github.com/toddnguyen47/splith1headers/internal/imagefix"
	"github.com/toddnguyen47/splith1headers/internal/splith1headers"
)

func main() {

	flagsRead := cmdlineflags.ParseCommandLineFlags()

	switch flagsRead.Command {
	case cmdlineflags.CommandSplit:
		splitHeaderStruct := splith1headers.NewSplitStruct()
		splitHeaderStruct.Split(flagsRead.BookPath, flagsRead.SplitImages)
	case cmdlineflags.CommandImageFix:
		imagefix.FixImages(flagsRead.BookPath)
	default:
		panic("Invalid command. Valid commands are: " + strings.Join(cmdlineflags.ValidCommmands, ", "))
	}
}
