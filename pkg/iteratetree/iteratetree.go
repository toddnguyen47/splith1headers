package iteratetree

import (
	"strings"

	"github.com/beevik/etree"
)

func IterateToFindTag(root *etree.Element, finalTag string, tags []string,
	curIndex int) *etree.Element {
	if root == nil || !strings.EqualFold(root.Tag, tags[curIndex]) {
		return nil
	}
	if curIndex == len(tags)-1 {
		// We found it!
		return root
	}

	var elem *etree.Element
	childElems := root.ChildElements()
	for _, childElem := range childElems {
		elem = IterateToFindTag(childElem, finalTag, tags, curIndex+1)
		if elem != nil && strings.EqualFold(finalTag, elem.Tag) {
			break
		}
	}

	return elem
}
