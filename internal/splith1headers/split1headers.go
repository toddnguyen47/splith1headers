package splith1headers

import (
	"container/list"
	"fmt"
	"strings"

	"github.com/beevik/etree"
	"github.com/toddnguyen47/splith1headers/pkg/iteratetree"
)

type splitStruct struct {
	elems      [][]*etree.Element
	index      int
	reverseMap map[string]int
	stack      *list.List
	header     *etree.Element
}

func NewSplitStruct() splitStruct {
	return splitStruct{
		elems:      make([][]*etree.Element, 0),
		index:      0,
		reverseMap: make(map[string]int),
		stack:      list.New(),
		header:     etree.NewElement("head"),
	}
}

// Split - split files into separate h1 headers.
//
// PARAMETERS
// inputFile - the path of the inputFile
// splitImages - true if you want images to be in its own xhtml file, false otherwise
func (s *splitStruct) Split(inputFile string, splitImages bool) {
	fmt.Printf("INIT Split() with inputFile: `%s`, splitImages: `%t`\n", inputFile, splitImages)
	doc := etree.NewDocument()
	err := doc.ReadFromFile(inputFile)
	if err != nil {
		panic(fmt.Sprintf("ERROR reading from file: %v", err))
	}

	s.header = iteratetree.IterateToFindTag(&doc.Element, "head", []string{"", "html", "head"}, 0)
	bodyElem := iteratetree.IterateToFindTag(&doc.Element, "body", []string{"", "html", "body"}, 0)

	// Recursively parse elements
	for _, childElem := range bodyElem.ChildElements() {
		s.parseTree(childElem, splitImages)
	}
	s.writeToFiles()

	fmt.Printf("END Split() with inputFile: `%s`, splitImages: `%t`\n", inputFile, splitImages)
}

func (s *splitStruct) parseTree(root *etree.Element, splitImages bool) {

	elemToAppend := root
	elemTag := strings.ToLower(elemToAppend.Tag)

	// If last tag was an img, make a new page
	if splitImages && s.stack.Len() > 0 {
		tailElem := s.stack.Back()
		val := tailElem.Value.(string)
		if strings.EqualFold(val, "img") {
			s.makeNewXmlList()
		}
	}

	if _, isHeader := constants.primaryHeadersMap[elemTag]; isHeader {
		// If h1/h2/h3, advance to the next file
		s.makeNewXmlList()
	} else if strings.EqualFold("div", elemTag) {
		elemToAppend = s.handleImages(elemToAppend, splitImages)
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

	if s.index >= len(s.elems) {
		s.makeNewXmlList()
	}
	s.elems[s.index] = append(s.elems[s.index], elemToAppend)

	// Recursively parse elements
	// for _, childElem := range root.ChildElements() {
	// 	s.parseTree(childElem, splitImages)
	// }
}

func (s *splitStruct) makeNewXmlList() {
	s.elems = append(s.elems, make([]*etree.Element, 0))
	s.index = len(s.elems) - 1
}

func (s *splitStruct) getFileName(index int) string {
	fileName := fmt.Sprintf("%s%02d.xhtml", constants.baseFileName, index)
	return fileName
}
