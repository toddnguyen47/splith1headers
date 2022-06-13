package splith1headers

import (
	"github.com/beevik/etree"
	"github.com/toddnguyen47/splith1headers/pkg/iteratetree"
)

func (s *splitStruct) handleImages(elemToAppend *etree.Element, splitImages bool) *etree.Element {
	imageElem := iteratetree.IterateToFindTag(elemToAppend, "image", []string{"div", "svg", "image"}, 0)
	if imageElem != nil {
		newElem := etree.NewElement("img")
		newAttr := etree.Attr{
			Key:   "src",
			Value: imageElem.SelectAttrValue(constants.href, ""),
		}
		newElem.Attr = append(newElem.Attr, newAttr)
		elemToAppend = newElem

		if splitImages {
			s.makeNewXmlList()
		}
	}
	return elemToAppend
}
