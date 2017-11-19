package palitra

type imageDim struct {
	horizontalPixels int
	verticalPixels   int
}

type rgbT [3]int

type colorMapT struct {
	key   string
	value int
}

// pair is a simple key value struct, which links a color with its hex representation
type pair struct {
	color string
	hex   string
}

// ColorPercentageT is a struct which contains the approximated CSS3 color `Color` and its
// spread in the image, expressed as `Percentage`%
type ColorPercentageT struct {
	Color      string `json:"color"`
	Percentage string `json:"percentage"`
}
