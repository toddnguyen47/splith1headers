package splith1headers

import (
	"fmt"
	"os"
	"path"

	"github.com/beevik/etree"
)

func (s *splitStruct) writeToFiles() {
	fmt.Printf("Map: %v\n", s.reverseMap)
	err := os.RemoveAll(constants.outputFolder)
	if err != nil {
		panic(fmt.Sprintf("ERROR while removing folder: %s. Error: %v", constants.outputFolder, err))
	}
	err = os.Mkdir(constants.outputFolder, os.ModeDir)
	if err != nil {
		panic(fmt.Sprintf("ERROR while making folder: %s. Error: %v", constants.outputFolder, err))
	}

	for i, xmlInFile := range s.elems {
		fileName := s.getFileName(i)
		fullFileName := path.Join(constants.outputFolder, fileName)
		rootElem := etree.NewElement("div")
		rootElem.CreateAttr("xmlns", "http://www.w3.org/1999/xhtml")
		for _, xmlElem := range xmlInFile {
			rootElem.AddChild(xmlElem)
		}
		newDoc := etree.NewDocument()
		newDoc.SetRoot(rootElem)
		newDoc.Indent(2)

		fmt.Println("Writing to file " + fullFileName)
		err := newDoc.WriteToFile(fullFileName)
		if err != nil {
			panic(fmt.Sprintf("ERROR writing to file %s. Error: %v", fullFileName, err))
		}
	}
}
