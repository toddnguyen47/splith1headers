package splith1headers

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/beevik/etree"
	"github.com/toddnguyen47/splith1headers/pkg/iteratetree"
)

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
		panic(fmt.Sprintf("ERROR reading from file: %v", err))
	}

	bodyElem := iteratetree.IterateToFindTag(&doc.Element, "body", []string{"", "html", "body"}, 0)

	// Recursively parse elements
	for _, childElem := range bodyElem.ChildElements() {
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
		imageElem := iteratetree.IterateToFindTag(elemToAppend, "image", []string{"div", "svg", "image"}, 0)
		if imageElem != nil {
			newElem := etree.NewElement("img")
			newAttr := etree.Attr{
				Key:   "src",
				Value: imageElem.SelectAttrValue(constants.href, ""),
			}
			newElem.Attr = append(newElem.Attr, newAttr)
			elemToAppend = newElem
		}
	} else if strings.EqualFold("p", elemToAppend.Tag) {
		for _, childElem := range elemToAppend.ChildElements() {
			if strings.EqualFold(childElem.Tag, "sup") {
				supId := childElem.SelectAttr("id")
				if supId != nil {
					s.reverseMap["#"+supId.Value] = s.index
				}

				aElem := iteratetree.IterateToFindTag(childElem, "a", []string{"sup", "a"}, 0)
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
				aElem := iteratetree.IterateToFindTag(childElem, "a", []string{"li", "span", "a"}, 0)
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

func (s *splitStruct) getFileName(index int) string {
	fileName := fmt.Sprintf("%s%02d.xhtml", constants.baseFileName, index)

	// If last file and file has notes
	if index == len(s.elems)-1 && s.hasNotesSection {
		fileName = fmt.Sprintf("%s.xhtml", constants.notes)
	}
	return fileName
}
