package imagefix

import (
	"fmt"
	"os"
	"path"

	"github.com/beevik/etree"
)

const (
	outputFolder = "output"
)

func Output(header *etree.Element, body *etree.Element, filename string, firstOutput bool) {
	if firstOutput {
		setupOutputFolder()
	}

	rootElem := setupRootElem(header, body)
	doc := setupNewDoc(rootElem)

	writeToFile(filename, doc)
}

func setupRootElem(header *etree.Element, body *etree.Element) *etree.Element {
	rootElem := etree.NewElement("html")
	rootElem.CreateAttr("xmlns", "http://www.w3.org/1999/xhtml")
	rootElem.AddChild(header)
	rootElem.AddChild(body)
	return rootElem
}

func setupNewDoc(rootElem *etree.Element) *etree.Document {
	newDoc := etree.NewDocument()
	newDoc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	newDoc.SetRoot(rootElem)
	newDoc.Indent(2)
	return newDoc
}

func setupOutputFolder() {
	err := os.RemoveAll(outputFolder)
	if err != nil {
		panic(fmt.Sprintf("ERROR while removing folder: %s. Error: %v", outputFolder, err))
	}

	err = os.Mkdir(outputFolder, os.ModeDir)
	if err != nil {
		panic(fmt.Sprintf("ERROR while making folder: %s. Error: %v", outputFolder, err))
	}
}

func writeToFile(filename string, newDoc *etree.Document) {
	fullFileName := path.Join(outputFolder, filename)
	fmt.Println("Writing to file " + fullFileName)
	err := newDoc.WriteToFile(fullFileName)
	if err != nil {
		panic(fmt.Sprintf("ERROR writing to file %s. Error: %v", fullFileName, err))
	}
}
