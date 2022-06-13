package splith1headers

import (
	"fmt"
	"strings"

	"github.com/beevik/etree"
	"github.com/toddnguyen47/splith1headers/pkg/iteratetree"
)

// handleCiteNote - Fix `cite_note` <sup> superscript text
func (s *splitStruct) handleCiteNote(elemToAppend *etree.Element) {
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
}
