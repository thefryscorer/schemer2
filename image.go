package main

import (
	"errors"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"os"
	"time"
)

func loadImage(filepath string) image.Image {
	infile, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer infile.Close()

	src, _, err := image.Decode(infile)
	if err != nil {
		log.Fatal(err)
	}
	return src
}

func abs(n int) int {
	if n >= 0 {
		return n
	}
	return -n
}

func colorDifference(col1 color.Color, col2 color.Color, threshold int) bool {
	c1 := col1.(color.NRGBA)
	c2 := col2.(color.NRGBA)

	rDiff := abs(int(c1.R) - int(c2.R))
	gDiff := abs(int(c1.G) - int(c2.G))
	bDiff := abs(int(c1.B) - int(c2.B))

	total := rDiff + gDiff + bDiff
	return total >= threshold
}

func getDistinctColors(colors []color.Color, threshold int, minBrightness, maxBrightness int) []color.Color {
	distinctColors := make([]color.Color, 0)
	for _, c := range colors {
		same := false
		if !colorDifference(c, color.NRGBAModel.Convert(color.Black), minBrightness*3) {
			continue
		}
		if !colorDifference(c, color.NRGBAModel.Convert(color.White), (255-maxBrightness)*3) {
			continue
		}
		for _, k := range distinctColors {
			if !colorDifference(c, k, threshold) {
				same = true
				break
			}
		}
		if !same {
			distinctColors = append(distinctColors, c)
		}
	}
	return distinctColors
}

func colorsFromImage(filename string) ([]color.Color, error) {
	// Load the image and create array of colors
	fuzzyness := 5
	img := loadImage(filename)
	w, h := img.Bounds().Max.X, img.Bounds().Max.Y
	colors := make([]color.Color, 0, w*h)
	for x := 0; x < w; x += fuzzyness {
		for y := 0; y < h; y += fuzzyness {
			col := color.NRGBAModel.Convert(img.At(x, y))
			colors = append(colors, col)
		}
	}
	// Get the distinct colors from the array by comparing differences with a threshold
	distinctColors := getDistinctColors(colors, *threshold, *minBrightness, *maxBrightness)

	// Ensure there are 16 colors
	count := 0
	for len(distinctColors) < 16 {
		count++
		distinctColors = append(distinctColors, getDistinctColors(colors, *threshold-count, *minBrightness, *maxBrightness)...)
		if count == *threshold {
			return nil, errors.New("Could not get colors from image with settings specified. Aborting.\n")
		}
	}

	if len(distinctColors) > 16 {
		distinctColors = distinctColors[:16]
	}

	return distinctColors, nil
}

func imageFromColors(colors []color.Color, w int, h int) image.Image {
	rand.Seed(time.Now().UnixNano())
	switch rand.Intn(4) {
	case 0:
		// Circles
		switch rand.Intn(2) {
		case 0:
			return Circles(colors, w, h, false)
		case 1:
			return Circles(colors, w, h, true)
		}
	case 1:
		// Rays
		switch rand.Intn(2) {
		case 0:
			return Rays(colors, w, h, true, rand.Intn(w/24))
		case 1:
			return Rays(colors, w, h, false, rand.Intn(w/24))
		}
	case 2:
		// Horizontal Lines
		switch rand.Intn(2) {
		case 0:
			return HorizontalLines(colors, w, h, false)
		case 1:
			return HorizontalLines(colors, w, h, true)
		}
	case 3:
		// Vertical Lines
		switch rand.Intn(4) {
		case 0:
			return VerticalLines(colors, w, h, false, false)
		case 1:
			return VerticalLines(colors, w, h, true, false)
		case 2:
			return VerticalLines(colors, w, h, false, true)
		case 3:
			return VerticalLines(colors, w, h, true, true)
		}
	}
	return nil
}

type Circle struct {
	col  color.Color
	x, y int
	size int
}

func Circles(colors []color.Color, w int, h int, filled bool) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))

	circles := make([]Circle, 0)

	for _, c := range colors {
		circle := Circle{c, rand.Intn(w), rand.Intn(h), rand.Intn(w / 2)}
		circles = append(circles, circle)
	}

	bg := colors[0]
	border := rand.Intn(w / 24)

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			img.Set(x, y, bg)
			for _, c := range circles {
				a := float64((x - c.x) * (x - c.x))
				b := float64((y - c.y) * (y - c.y))

				if filled {
					if int(math.Sqrt(a+b)) < c.size {
						img.Set(x, y, c.col)
					}
				} else {
					if int(math.Sqrt(a+b)) < c.size && int(math.Sqrt(a+b)) > (c.size-border) {
						img.Set(x, y, c.col)
					}
				}
			}
		}
	}
	return img
}

type Stripe struct {
	col   color.Color
	x, y  int // Middle point
	angle int // 0-180
}

func Rays(colors []color.Color, w int, h int, centered bool, margin int) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))

	stripes := make([]Stripe, 0)

	for _, c := range colors {
		var stripe Stripe
		if centered {
			stripe = Stripe{c, w / 2, h / 2, rand.Intn(180)}
		} else {
			stripe = Stripe{c, rand.Intn(w), rand.Intn(h), rand.Intn(180)}
		}
		stripes = append(stripes, stripe)
	}

	bg := colors[0]

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			img.Set(x, y, bg)
			for _, s := range stripes {
				deltaX := float64(x - s.x)
				deltaY := float64(y - s.y)
				angle := math.Atan(deltaY/deltaX) * 180 / math.Pi
				if int(math.Abs(float64(int(angle)-s.angle))) < margin {
					img.Set(x, y, s.col)
				}
			}
		}
	}
	return img
}

type VerticalLine struct {
	col color.Color
	x   int
	w   int
}

type HorizontalLine struct {
	col color.Color
	y   int
	h   int
}

func VerticalLines(colors []color.Color, w int, h int, evenlySpaced bool, evenWidth bool) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))

	lines := make([]VerticalLine, 0)

	var width int
	width = rand.Intn(w / 16)

	x_index := rand.Intn(w / 2)

	var spacing int
	spacing = rand.Intn(w / 32)

	for _, c := range colors {
		if !evenWidth {
			width = rand.Intn(w / 16)
		}
		if !evenlySpaced {
			spacing = rand.Intn(w / 32)
		}
		x_index += spacing
		lines = append(lines, VerticalLine{c, x_index, width})
		x_index += width
	}

	bg := colors[0]

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			img.Set(x, y, bg)
			for _, l := range lines {
				if x >= l.x && x < l.x+l.w {
					img.Set(x, y, l.col)
				}
			}
		}
	}

	return img

}

func HorizontalLines(colors []color.Color, w int, h int, evenHeight bool) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))

	lines := make([]HorizontalLine, 0)

	var height int
	if evenHeight {
		height = rand.Intn(h / 16)
	}

	for _, c := range colors {
		if !evenHeight {
			height = rand.Intn(h / 16)
		}
		lines = append(lines, HorizontalLine{c, rand.Intn(h), height})
	}

	bg := colors[0]

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			img.Set(x, y, bg)
			for _, l := range lines {
				if y >= l.y && y < l.y+l.h {
					img.Set(x, y, l.col)
				}
			}
		}
	}

	return img

}
