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
	notes        string
	href         string
}{
	baseFileName: "Body",
	outputFolder: "output",
	notes:        "Notes",
	href:         "href",
}

type splitStruct struct {
	elems           [][]*etree.Element
	index           int
	hasNotesSection bool
	reverseMap      map[string]int
}

func NewSplitStruct() splitStruct {
	return splitStruct{
		elems:           make([][]*etree.Element, 0),
		index:           0,
		hasNotesSection: false,
		reverseMap:      make(map[string]int),
	}
}

func (s *splitStruct) Split(inputFile string) {
	doc := etree.NewDocument()
	err := doc.ReadFromFile(inputFile)
	if err != nil {
		fmt.Println("ERROR reading from file")
	}

	// Recursively parse elements
	for _, childElem := range doc.Element.ChildElements() {
		s.parseTree(childElem)
	}
	s.writeToFiles()
}

func (s *splitStruct) parseTree(root *etree.Element) {

	elemToAppend := root
	// If h1, advance to the next file
	if strings.EqualFold(elemToAppend.Tag, "h1") {
		s.elems = append(s.elems, make([]*etree.Element, 0))
		s.index = len(s.elems) - 1
		// If notes
		attributes := elemToAppend.Attr
		for _, attribute := range attributes {
			if strings.EqualFold(attribute.Key, "id") && strings.EqualFold(attribute.Value, constants.notes) {
				s.hasNotesSection = true
			}
		}
	} else if strings.EqualFold("div", elemToAppend.Tag) {
		// Select image
		imageElem := elemToAppend.FindElement("//svg/image")
		newElem := etree.NewElement("img")
		newAttr := etree.Attr{
			Key:   "src",
			Value: imageElem.SelectAttrValue(constants.href, ""),
		}
		newElem.Attr = append(newElem.Attr, newAttr)
		elemToAppend = newElem
	} else if strings.EqualFold("p", elemToAppend.Tag) {
		for _, childElem := range elemToAppend.ChildElements() {
			if strings.EqualFold(childElem.Tag, "sup") {
				supId := childElem.SelectAttr("id")
				if supId != nil {
					s.reverseMap["#"+supId.Value] = s.index
				}

				aElem := childElem.FindElement("//sup/a")
				attr := aElem.SelectAttr(constants.href)
				if strings.Contains(attr.Value, "cite_note") {
					attr.Value = fmt.Sprintf("Notes.xhtml%s", attr.Value)
				}
				aElem.CreateAttr(constants.href, attr.Value)
			}
		}
	} else if strings.EqualFold("ol", elemToAppend.Tag) {
		// For each <li> element
		for _, childElem := range elemToAppend.ChildElements() {
			if strings.EqualFold("li", childElem.Tag) {
				aElem := childElem.FindElement("//li/span/a")
				attr := aElem.SelectAttr(constants.href)
				if index, ok := s.reverseMap[attr.Value]; ok {
					attr.Value = fmt.Sprintf("%s%s", s.getFileName(index), attr.Value)
					aElem.CreateAttr(constants.href, attr.Value)
				}
			}
		}
	}

	s.elems[s.index] = append(s.elems[s.index], elemToAppend)

	// Recursively parse elements
	// for _, childElem := range root.ChildElements() {
	// 	s.parseTree(childElem)
	// }
}

func (s *splitStruct) writeToFiles() {
	fmt.Printf("Map: %v\n", s.reverseMap)
	os.Mkdir(constants.outputFolder, os.ModeDir)
	for i, xmlInFile := range s.elems {
		fileName := s.getFileName(i)
		fullFileName := path.Join(constants.outputFolder, fileName)
		rootElem := etree.NewElement("div")
		for _, xmlElem := range xmlInFile {
			rootElem.AddChild(xmlElem)
		}
		newDoc := etree.NewDocument()
		newDoc.SetRoot(rootElem)
		newDoc.Indent(2)

		fmt.Println("Writing to file " + fullFileName)
		err := newDoc.WriteToFile(fullFileName)
		if err != nil {
			fmt.Println("ERROR writing to file " + fullFileName)
		}
	}
}

func (s *splitStruct) getFileName(index int) string {
	fileName := fmt.Sprintf("%s%02d.xhtml", constants.baseFileName, index)

	// If last file and file has notes
	if index == len(s.elems)-1 && s.hasNotesSection {
		fileName = fmt.Sprintf("%s.xhtml", constants.notes)
	}
	return fileName
}
