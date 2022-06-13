package splith1headers

import (
	"fmt"
	"strings"

	"github.com/beevik/etree"
	"github.com/toddnguyen47/splith1headers/pkg/iteratetree"
)

// handleLists - For each <li> element
func (s *splitStruct) handleLists(elemToAppend *etree.Element) {
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
