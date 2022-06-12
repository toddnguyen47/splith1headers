package main

import "github.com/toddnguyen47/splith1headers/internal/splith1headers"

func main() {
	inputFile := `Z:/DocumentsAndStuff/Desktop/Body.xhtml`

	splitHeaderStruct := splith1headers.NewSplitStruct()
	splitHeaderStruct.Split(inputFile)
}
