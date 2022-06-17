package main

import (
	"github.com/toddnguyen47/splith1headers/internal/cmdlineflags"
	"github.com/toddnguyen47/splith1headers/internal/splith1headers"
)

func main() {

	flagsRead := cmdlineflags.ParseCommandLineFlags()

	splitHeaderStruct := splith1headers.NewSplitStruct()
	splitHeaderStruct.Split(flagsRead.BookPath, flagsRead.SplitImages)
}
