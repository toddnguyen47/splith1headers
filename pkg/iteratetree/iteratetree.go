package iteratetree

import (
	"strings"

	"github.com/beevik/etree"
)

func IterateToFindTag(root *etree.Element, finalTag string, tags []string,
	curIndex int) *etree.Element {
	if curIndex >= len(tags) || strings.EqualFold(finalTag, tags[curIndex]) {
		return root
	}

	var elem *etree.Element
	for _, childElem := range root.ChildElements() {
		elem = IterateToFindTag(childElem, finalTag, tags, curIndex+1)
		if elem != nil {
			break
		}
	}

	return elem
}
