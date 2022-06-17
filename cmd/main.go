package main

import (
	"github.com/toddnguyen47/splith1headers/internal/cmdlineflags"
	"github.com/toddnguyen47/splith1headers/internal/splith1headers"
)

func main() {

	inputFile, splitImages := cmdlineflags.ParseCommandLineFlags()

	splitHeaderStruct := splith1headers.NewSplitStruct()
	splitHeaderStruct.Split(*inputFile, *splitImages)
}
