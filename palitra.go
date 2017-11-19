package palitra

import "image"

// getColorMap populates a map of approximated colors, given the color of a single
// pixel. To enhance performance, it does so in a concurrent fashion, using channels and
// an anonymous goroutine.
// It returns the color map and the total number of pixels of the image `img`.
func getColorMap(img image.Image) (map[string]int, int) {
	colorMap := make(map[string]int)
	dims := getImageDim(img)
	totalPixels := dims.horizontalPixels * dims.verticalPixels

	approxColorResponseChan := make(chan string, totalPixels)
	defer close(approxColorResponseChan)

	for i := 0; i < dims.horizontalPixels; i++ {
		for j := 0; j < dims.verticalPixels; j++ {
			pixel := img.At(i, j)
			rgbTuple := getRGBTuple(pixel)

			go func(rgb rgbT) {
				approxColorResponseChan <- getApproximatedCSSColor(rgb)
			}(rgbTuple)
		}
	}

	for k := 0; k < totalPixels; k++ {
		colorName := <-approxColorResponseChan
		_, exists := colorMap[colorName]

		if exists {
			colorMap[colorName]++
		} else {
			colorMap[colorName] = 1
		}
	}

	return colorMap, totalPixels
}

// GetPalette returns an array of `paletteLen` colors, given an image `img`. To enhance
// performance, a `resizeWidth` param is required.
func GetPalette(img image.Image, paletteLimit int, resizeWidth uint) []ColorPercentageT {
	resizedImg := resizeImage(img, resizeWidth)
	colorMap, totalPixels := getColorMap(resizedImg)

	palette := sortMap(colorMap)
	paletteWithPercentage := getPaletteWithPercentage(palette, paletteLimit, float64(totalPixels))
	return paletteWithPercentage
}
