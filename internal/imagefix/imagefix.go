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
	newBody := etree.NewElement("body")
	queue := list.New()

	for _, childElem := range bodyElem.ChildElements() {
		// Recursively parse elements
		parseTree(childElem, bodyElem, queue)
		newBody.AddChild(childElem)

		for queue.Len() > 0 {
			node := queue.Front()
			queue.Remove(node)
			newBody.AddChild(node.Value.(*etree.Element))
		}
	}

	_, filename := filepath.Split(inputFile)
	Output(header, newBody, filename, firstOutput)

	fmt.Printf("END FixImage() with inputFile: `%s`\n", inputFile)
}

func parseTree(root *etree.Element, parent *etree.Element, queue *list.List) {

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
		parent.RemoveChild(root)
	}

	// Recursively parse elements
	for _, childElem := range root.ChildElements() {
		parseTree(childElem, root, queue)
	}

	// If queue is not empty, add "hidden" to current element
	if queue.Len() > 0 && !strings.EqualFold("body", root.Tag) {
		root.CreateAttr("class", "hidden")
	}
}
