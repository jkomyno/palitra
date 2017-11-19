package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"os"
	"time"

	"github.com/jkomyno/palitra"
)

var (
	inputImageName = flag.String("img", "", "the name of the image file to get colors from")
	paletteLimit   = flag.Int("n", 3, "the number of the most dominant CSS3 colors to extract")
	resizeWidth    = flag.Uint("r", 75, "the resized width to process the image quicker")
)

// supportsStdin returns true if os.Stdin is interactive
func supportsStdin() bool {
	fileInfo, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return fileInfo.Mode()&(os.ModeCharDevice|os.ModeCharDevice) != 0
}

// showSpinner is a simple loader to display the user that palitra is processing
func showSpinner() {
	for { // while true
		for _, r := range `⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏` {
			fmt.Printf("\r%c %s", r, "Calculating colors")
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func timeTrack(start time.Time) {
	elapsed := time.Since(start)
	fmt.Printf("\nCompleted in %s ✨", elapsed)
}

func printPaletteList(palette []palitra.ColorPercentageT) {
	fmt.Printf("\r")
	for i, val := range palette {
		fmt.Printf("%d) %s: %s%s\n", i+1, val.Color, val.Percentage, "%")
	}
}

func executeCommand(img image.Image) {
	defer timeTrack(time.Now())
	palette := palitra.GetPalette(img, *paletteLimit, *resizeWidth)
	printPaletteList(palette)
}

func init() {
	flag.Parse()
}

func main() {
	// handle inputImageName

	if supportsStdin() && *inputImageName == "" {
		flag.Usage()
		fmt.Fprintln(os.Stderr, "Supports reading image from stdin")
		os.Exit(1)
	}

	var imageReader io.Reader
	imageReader = os.Stdin
	if *inputImageName != "" {
		f, err := os.Open(*inputImageName)
		if err != nil {
			log.Fatalf("Reading input file: %s", err)
		}

		defer f.Close()
		imageReader = f
	}

	if img, _, err := image.Decode(imageReader); err != nil {
		fmt.Fprintln(os.Stderr, "error parsing", err)
		os.Exit(1)
	} else {
		go showSpinner()
		executeCommand(img)
	}
}
