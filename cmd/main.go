package main

import "github.com/toddnguyen47/splith1headers/internal/splith1headers"

func main() {
	inputFile := `Z:/DocumentsAndStuff/Desktop/1.xhtml`

	splitHeaderStruct := splith1headers.NewSplitStruct()
	splitHeaderStruct.Split(inputFile)
}
