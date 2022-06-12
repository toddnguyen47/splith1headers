package splith1headers

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/beevik/etree"
)

var constants = struct {
	baseFileName string
	outputFolder string
}{
	baseFileName: "body",
	outputFolder: "output",
}

type splitStruct struct {
	elems [][]*etree.Element
	index int
}

func NewSplitStruct() splitStruct {
	return splitStruct{
		elems: make([][]*etree.Element, 0),
		index: 0,
	}
}

func (s *splitStruct) Split(inputFile string) {
	doc := etree.NewDocument()
	err := doc.ReadFromFile(inputFile)
	if err != nil {
		fmt.Println("ERROR reading from file")
	}
	s.parseTree(&doc.Element)
	// TODO: WRITE TO FILES
	s.writeToFiles()
}

func (s *splitStruct) parseTree(root *etree.Element) {
	if !isElemSkipped(root) {
		elemToAppend := root
		// If h1, advance to the next file
		if root.Tag == "h1" {
			s.elems = append(s.elems, make([]*etree.Element, 0))
			s.index = len(s.elems) - 1
		} else if root.Tag == "image" {
			newElem := etree.NewElement("img")
			newAttr := etree.Attr{
				Key:   "src",
				Value: root.SelectAttrValue("href", ""),
			}
			newElem.Attr = append(newElem.Attr, newAttr)
			elemToAppend = newElem
		}

		s.elems[s.index] = append(s.elems[s.index], elemToAppend)
	}

	// Recursively parse elements
	for _, childElem := range root.ChildElements() {
		s.parseTree(childElem)
	}
}

func (s *splitStruct) writeToFiles() {
	os.Mkdir(constants.outputFolder, os.ModeDir)
	for i, xmlInFile := range s.elems {
		fileName := fmt.Sprintf("%s%02d.xhtml", constants.baseFileName, i)
		fullFileName := path.Join(constants.outputFolder, fileName)
		rootElem := etree.NewElement("div")
		for _, xmlElem := range xmlInFile {
			rootElem.AddChild(xmlElem)
		}
		newDoc := etree.NewDocument()
		newDoc.SetRoot(rootElem)

		fmt.Println("Writing to file " + fullFileName)
		err := newDoc.WriteToFile(fullFileName)
		if err != nil {
			fmt.Println("ERROR writing to file " + fullFileName)
		}
	}
}

func isElemSkipped(root *etree.Element) bool {
	// Empty tags
	if len(root.Tag) <= 0 {
		return true
	}

	// Skip <svg> and <body>
	if strings.EqualFold(root.Tag, "svg") || strings.EqualFold(root.Tag, "body") {
		return true
	}

	// Skip <div class=svg_outer
	if strings.EqualFold("div", root.Tag) {
		attributes := root.Attr
		for _, attribute := range attributes {
			if strings.EqualFold(attribute.Key, "class") && strings.Contains(attribute.Value, "svg_outer") {
				return true
			}
		}
	}

	return false
}
