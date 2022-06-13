package splith1headers

import (
	"container/list"
	"fmt"
	"strings"

	"github.com/beevik/etree"
	"github.com/toddnguyen47/splith1headers/pkg/iteratetree"
)

type splitStruct struct {
	elems           [][]*etree.Element
	index           int
	hasNotesSection bool
	reverseMap      map[string]int
	stack           *list.List
}

func NewSplitStruct() splitStruct {
	return splitStruct{
		elems:           make([][]*etree.Element, 0),
		index:           0,
		hasNotesSection: false,
		reverseMap:      make(map[string]int),
		stack:           list.New(),
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
	elemTag := strings.ToLower(elemToAppend.Tag)

	// If last tag was an img, make a new page
	if s.stack.Len() > 0 {
		tailElem := s.stack.Back()
		val := tailElem.Value.(string)
		if strings.EqualFold(val, "img") {
			s.makeNewXmlList()
		}
	}

	// If h1, advance to the next file
	if _, isHeader := constants.primaryHeadersMap[elemTag]; isHeader {
		s.handleHeaders(elemToAppend)
	} else if strings.EqualFold("div", elemTag) {
		elemToAppend = s.handleImages(elemToAppend)
	} else if strings.EqualFold("p", elemTag) {
		s.handleCiteNote(elemToAppend)
	} else if strings.EqualFold("ol", elemTag) {
		s.handleLists(elemToAppend)
	}

	// Pop last element and push this element
	if tailElem := s.stack.Back(); tailElem != nil {
		s.stack.Remove(tailElem)
	}
	s.stack.PushBack(elemToAppend.Tag)

	s.elems[s.index] = append(s.elems[s.index], elemToAppend)

	// Recursively parse elements
	// for _, childElem := range root.ChildElements() {
	// 	s.parseTree(childElem)
	// }
}

func (s *splitStruct) makeNewXmlList() {
	s.elems = append(s.elems, make([]*etree.Element, 0))
	s.index = len(s.elems) - 1
}

func (s *splitStruct) getFileName(index int) string {
	fileName := fmt.Sprintf("%s%02d.xhtml", constants.baseFileName, index)

	// If last file and file has notes
	if index == len(s.elems)-1 && s.hasNotesSection {
		fileName = fmt.Sprintf("%s.xhtml", constants.notes)
	}
	return fileName
}
