package splith1headers

import (
	"strings"

	"github.com/beevik/etree"
)

func (s *splitStruct) handleHeaders(elemToAppend *etree.Element) {
	s.elems = append(s.elems, make([]*etree.Element, 0))
	s.index = len(s.elems) - 1

	attributes := elemToAppend.Attr
	for _, attribute := range attributes {
		if strings.EqualFold(attribute.Key, "id") && strings.EqualFold(attribute.Value, constants.notes) {
			s.hasNotesSection = true
		}
	}
}
