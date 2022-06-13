package splith1headers

var constants = struct {
	baseFileName      string
	outputFolder      string
	notes             string
	href              string
	maxElemsPerFile   int
	primaryHeadersMap map[string]struct{}
}{
	baseFileName:    "Body",
	outputFolder:    "output",
	notes:           "Notes",
	href:            "href",
	maxElemsPerFile: 500,
	primaryHeadersMap: map[string]struct{}{
		"h1": {},
		"h2": {},
		"h3": {},
	},
}
