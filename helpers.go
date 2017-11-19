package palitra

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"sort"

	webColors "github.com/jyotiska/go-webcolors"
	"github.com/nfnt/resize"
)

// l2Norm returns the Euclidean distance between two colors
func l2Norm(currentColor rgbT, cssColor []int) int {
	r := cssColor[0] - currentColor[0]
	g := cssColor[1] - currentColor[1]
	b := cssColor[2] - currentColor[2]

	return int(math.Sqrt(float64(r*r + g*g + b*b)))
}

// getApproximatedCSSColor returns the nearest approximated CSS3 color of the color
// passed as parameter
func getApproximatedCSSColor(currentColor rgbT) string {
	approximateColors := make(map[int]string)

	for _, c := range cssColorPair {
		rgb := webColors.HexToRGB(c.hex)
		approximateColors[l2Norm(currentColor, rgb)] = c.color
	}

	keys := make([]int, 0, len(approximateColors))
	for key := range approximateColors {
		keys = append(keys, key)
	}
	sort.Ints(keys)

	return approximateColors[keys[0]]
}

// min returs the minimum value between `a` and `b`
func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

// resizeImage resizes an image adapting the aspect ration to a specified width
func resizeImage(img image.Image, width uint) image.Image {
	// height set to 0 -> adapts aspect ratio to specified width
	return resize.Resize(width, 0, img, resize.Lanczos3)
}

// getImageDim returns the horizontal and vertical dimensions of an image, in pixels
func getImageDim(img image.Image) *imageDim {
	maxBound := img.Bounds().Max
	return &imageDim{
		horizontalPixels: maxBound.X,
		verticalPixels:   maxBound.Y,
	}
}

// getRGBTuple, given a single pixel, retuns its RGB color
func getRGBTuple(pixel color.Color) rgbT {
	red, green, blue, _ := pixel.RGBA()
	rgbTuple := rgbT{int(red / 255), int(green / 255), int(blue / 255)}
	return rgbTuple
}

// sortMap sorts `colorMap` by its value property
func sortMap(colorMap map[string]int) []colorMapT {
	var sortedColorMap []colorMapT

	for key, value := range colorMap {
		sortedColorMap = append(sortedColorMap, colorMapT{
			key,
			value,
		})
	}

	sort.Slice(sortedColorMap, func(i, j int) bool {
		return sortedColorMap[i].value > sortedColorMap[j].value
	})

	return sortedColorMap
}

// getPaletteWithPercentage returns a maximum of `paletteLimit` combinations of approximated color
// and percentage. The percentage represents how much an approximation of that color appears in
// the original image, and for convenience is expressed as string.
func getPaletteWithPercentage(palette []colorMapT, paletteLimit int, totalPixels float64) []ColorPercentageT {
	size := min(paletteLimit, len(palette))
	paletteWithPercentage := make([]ColorPercentageT, size)
	var currPalette colorMapT

	for i := 0; i < size; i++ {
		currPalette = palette[i]
		paletteWithPercentage[i] = ColorPercentageT{
			Color:      currPalette.key,
			Percentage: fmt.Sprintf("%.2f", (float64(currPalette.value)/totalPixels)*100),
		}
	}

	return paletteWithPercentage
}
