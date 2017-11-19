package palitra_test

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"testing"

	"github.com/jkomyno/palitra"
	"github.com/stretchr/testify/assert"
)

var testImg image.Image

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	imageReader, err := os.Open("./images/venice.jpg")
	handleError(err)
	defer imageReader.Close()

	testImg, _, err = image.Decode(imageReader)
	handleError(err)
}

func TestGetPalette(t *testing.T) {
	colorMap1 := palitra.GetPalette(testImg, 3, 75)
	expectedMap1 := []palitra.ColorPercentageT{{
		Color:      "lightpink",
		Percentage: "10.48",
	}, {
		Color:      "lightslategrey",
		Percentage: "9.97",
	}, {
		Color:      "dimgrey",
		Percentage: "9.61",
	}}
	for i := range colorMap1 {
		assert.Equal(t, expectedMap1[i], colorMap1[i])
	}

	colorMap2 := palitra.GetPalette(testImg, 4, 75)
	expectedMap2 := append(expectedMap1, palitra.ColorPercentageT{
		Color:      "silver",
		Percentage: "9.36",
	})

	for i := range colorMap2 {
		assert.Equal(t, expectedMap2[i], colorMap2[i])
	}

	colorMap3 := palitra.GetPalette(testImg, 5, 75)
	expectedMap3 := append(expectedMap2, palitra.ColorPercentageT{
		Color:      "darkslategrey",
		Percentage: "8.79",
	})

	for i := range colorMap3 {
		assert.Equal(t, expectedMap3[i], colorMap3[i])
	}
}

// Benchmark nanoid generator
func BenchmarkGetPalette(b *testing.B) {

	for n := 0; n < b.N; n++ {
		palitra.GetPalette(testImg, 3, 75)
	}
}
