package splith1headers

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/beevik/etree"
)

func (s *splitStruct) writeToFiles() {
	fmt.Printf("Map: %v\n", s.reverseMap)
	setupOutputFolder()

	for i, xmlInFile := range s.elems {
		s.writePerFile(xmlInFile, i)
	}
}

func (s *splitStruct) writePerFile(xmlInFile []*etree.Element, index int) {
	isNotesSection, bodyElem := s.addToBodyElem(xmlInFile)
	rootElem := s.setupRootElem(bodyElem)
	newDoc := s.setupNewDoc(rootElem)
	s.writeToFile(index, isNotesSection, newDoc)
}

func (s *splitStruct) writeToFile(index int, isNotesSection bool, newDoc *etree.Document) {
	fileName := s.getFileName(index)
	fullFileName := path.Join(constants.outputFolder, fileName)
	if isNotesSection {
		fullFileName = path.Join(constants.outputFolder, constants.notes+".xhtml")
	}
	fmt.Println("Writing to file " + fullFileName)
	err := newDoc.WriteToFile(fullFileName)
	if err != nil {
		panic(fmt.Sprintf("ERROR writing to file %s. Error: %v", fullFileName, err))
	}
}

func (*splitStruct) setupNewDoc(rootElem *etree.Element) *etree.Document {
	newDoc := etree.NewDocument()
	newDoc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	newDoc.SetRoot(rootElem)
	newDoc.Indent(2)
	return newDoc
}

func (s *splitStruct) setupRootElem(bodyElem *etree.Element) *etree.Element {
	rootElem := etree.NewElement("html")
	rootElem.CreateAttr("xmlns", "http://www.w3.org/1999/xhtml")
	rootElem.AddChild(s.header)
	rootElem.AddChild(bodyElem)
	return rootElem
}

func (*splitStruct) addToBodyElem(xmlInFile []*etree.Element) (bool, *etree.Element) {
	isNotesSection := false
	bodyElem := etree.NewElement("body")
	for _, xmlElem := range xmlInFile {
		bodyElem.AddChild(xmlElem)
		if strings.EqualFold(xmlElem.Tag, "ol") && !isNotesSection {
			notesFound := xmlElem.FindElement("//ol/li/span[@class='mw-cite-backlink']")
			if notesFound != nil {
				isNotesSection = true
			}
		}
	}
	return isNotesSection, bodyElem
}

func setupOutputFolder() {
	err := os.RemoveAll(constants.outputFolder)
	if err != nil {
		panic(fmt.Sprintf("ERROR while removing folder: %s. Error: %v", constants.outputFolder, err))
	}
	err = os.Mkdir(constants.outputFolder, os.ModeDir)
	if err != nil {
		panic(fmt.Sprintf("ERROR while making folder: %s. Error: %v", constants.outputFolder, err))
	}
}
