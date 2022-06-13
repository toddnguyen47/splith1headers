package main

import "github.com/toddnguyen47/splith1headers/internal/splith1headers"

func main() {
	inputFile := "C:/Common/Desktop/Body.xhtml"

	splitHeaderStruct := splith1headers.NewSplitStruct()
	splitImages := true
	splitHeaderStruct.Split(inputFile, splitImages)
}
