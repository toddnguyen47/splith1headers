package imagefix

import (
	"container/list"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/beevik/etree"
	"github.com/toddnguyen47/splith1headers/pkg/iteratetree"
)

func FixImages(directory string) {
	files, err := os.ReadDir(directory)
	if err != nil {
		panic("Cannot read directory: " + directory)
	}

	firstFile := true
	for _, file := range files {
		fullpath := filepath.Join(directory, file.Name())
		FixImage(fullpath, firstFile)
		firstFile = false
	}
}

// FixImage - fix <image> link into <img> link
// For `firstOutput`, pass in True if it's the first file, false otherwise. This is needed
// for directory cleanup
func FixImage(inputFile string, firstOutput bool) {
	fmt.Printf("INIT FixImage() with inputFile: `%s`\n", inputFile)
	doc := etree.NewDocument()
	err := doc.ReadFromFile(inputFile)
	if err != nil {
		panic(fmt.Sprintf("ERROR reading from file: %v", err))
	}

	header := iteratetree.IterateToFindTag(&doc.Element, "head", []string{"", "html", "head"}, 0)
	bodyElem := iteratetree.IterateToFindTag(&doc.Element, "body", []string{"", "html", "body"}, 0)

	// Recursively parse elements
	newBody := etree.NewElement("body")
	queue := list.New()
	for _, childElem := range bodyElem.ChildElements() {
		parentsQueue := list.New()
		parseTree(childElem, parentsQueue, queue)

		frontNode := parentsQueue.Front()
		len1 := parentsQueue.Len()

		var node *list.Element
		for i := 0; i < len1-1; i++ {
			node = parentsQueue.Front()
			elem := node.Value.(*etree.Element)
			next := node.Next()
			elem.AddChild(next.Value.(*etree.Element))
			node = next
		}

		newBody.AddChild(frontNode.Value.(*etree.Element))

		for elem := queue.Front(); elem != nil; elem = elem.Next() {
			elemTree := elem.Value.(*etree.Element)
			newBody.AddChild(elemTree)
		}
	}

	_, filename := filepath.Split(inputFile)
	Output(header, newBody, filename, firstOutput)

	fmt.Printf("END FixImage() with inputFile: `%s`\n", inputFile)
}

func parseTree(root *etree.Element, parentsQueue *list.List, queue *list.List) {

	// If element is an image
	if strings.EqualFold("image", root.Tag) {
		// Create new element
		newElem := etree.NewElement("img")
		newElem.Attr = append(newElem.Attr, etree.Attr{
			Key:   "src",
			Value: root.SelectAttrValue("href", ""),
		})
		queue.PushBack(newElem)

		// Remove <image> element from parent
		lastParent := parentsQueue.Back().Value.(*etree.Element)
		lastParent.RemoveChild(root)

		repopulateQueue(parentsQueue)
	}

	parentsQueue.PushBack(root)
	// Recursively parse elements
	for _, childElem := range root.ChildElements() {
		parseTree(childElem, parentsQueue, queue)
	}
}

func repopulateQueue(parentsQueue *list.List) {
	// For every parent, remove all attributes
	newQueue := list.New()

	for parentsQueue.Len() > 0 {
		node := parentsQueue.Front()
		parentsQueue.Remove(node)

		elem := node.Value.(*etree.Element)
		elem.CreateAttr("class", "hidden")
		elem.Child = make([]etree.Token, 0)
		newQueue.PushBack(elem)
	}

	// Pop back in!
	for newQueue.Len() > 0 {
		node := newQueue.Front()
		newQueue.Remove(node)
		parentsQueue.PushBack(node.Value)
	}
}
