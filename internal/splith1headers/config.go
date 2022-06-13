package splith1headers

var constants = struct {
	baseFileName      string
	outputFolder      string
	notes             string
	href              string
	// Set `maxElemsPerFile` to -1 to have no limit
	maxElemsPerFile   int
	primaryHeadersMap map[string]struct{}
}{
	baseFileName:    "Body",
	outputFolder:    "output",
	notes:           "Notes",
	href:            "href",
	maxElemsPerFile: -1,
	primaryHeadersMap: map[string]struct{}{
		"h1": {},
		"h2": {},
		"h3": {},
	},
}
