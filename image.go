package main

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"time"
)

var imageOutTypes = [...]string{"random", "circles", "rays", "stripes"}

const m = 1<<16 - 1

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

func randMinMax(min int, max int) int {
	if min == max {
		return min
	}
	return rand.Intn(max-min) + min
}

func multiplyAlpha(c1 color.Color, c2 color.Color) color.Color {
	r1, g1, b1, a1 := color.RGBAModel.Convert(c1).RGBA()
	r2, g2, b2, a2 := color.RGBAModel.Convert(c2).RGBA()

	red1, green1, blue1, alpha1 := float64(r1), float64(g1), float64(b1), float64(a1)
	red2, green2, blue2, _ := float64(r2), float64(g2), float64(b2), float64(a2)

	if alpha1 == 255 {
		return c1
	}

	r := (red1 * alpha1 / 255) + red2*(1-alpha1/255)
	g := (green1 * alpha1 / 255) + green2*(1-alpha1/255)
	b := (blue1 * alpha1 / 255) + blue2*(1-alpha1/255)

	result := color.NRGBA{uint8(r), uint8(g), uint8(b), 255}
	return result
}

func randBool() bool {
	return rand.Intn(2) == 0
}

func capToMax(n, max int) int {
	if n > max {
		return max
	}
	return n
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

func imageFromColors(colors []color.Color, w int, h int) (image.Image, error) {
	rand.Seed(time.Now().UnixNano())
	switch *imageOutType {
	case "random":
		return randomImage(colors, w, h), nil
	case "circles":
		return Circles(colors, w, h, *circlesSize, *circlesSizeVariance, *circlesOverlap, *circlesDrawLargestToSmallest, *circlesFilled, *circlesBorderSize, *circlesBlur, *circlesOpacity), nil
	case "rays":
		return Rays(colors, w, h, *raysSize, *raysSizeVariance, *raysDistributeEvenly, *raysCentered, *raysDrawLargestToSmallest), nil
	case "stripes":
		return Lines(colors, w, h, *stripesSize, *stripesSizeVariance, *stripesHorizontal, *stripesEvenSpacing, *stripesSpacing, *stripesOffset), nil

	}
	return nil, errors.New("Unrecognised ouput image type: " + *imageOutType + "\n")
}

type Circle struct {
	col  color.Color
	x, y int
	size int
}

// For sorting circles by size
type circleBySize []Circle

func (a circleBySize) Len() int           { return len(a) }
func (a circleBySize) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a circleBySize) Less(i, j int) bool { return a[i].size < a[j].size }

func Circles(colors []color.Color, w int, h int, size int, sizevar int, overlap bool, large2small bool, filled bool, bordersize int, blur bool, opacity int) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))

	circles := make([]Circle, 0)
	bg := colors[0]

	for _, c := range colors {
		// Do not create circle with background color
		if c == bg {
			continue
		}
		circle := Circle{c, rand.Intn(w), rand.Intn(h), randMinMax(size-sizevar, size+sizevar)}
		circles = append(circles, circle)
	}

	if large2small {
		sort.Sort(sort.Reverse(circleBySize(circles)))
	}

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			img.Set(x, y, bg)
		}
	}
	for _, c := range circles {
		circleimg := image.NewNRGBA(image.Rect(0, 0, c.size*2, c.size*2))
		for x := 0; x < c.size*2; x++ {
			for y := 0; y < c.size*2; y++ {

				a := float64((x - c.size) * (x - c.size))
				b := float64((y - c.size) * (y - c.size))

				dist := int(math.Sqrt(a + b))

				var col color.Color
				if blur {
					alph := 255 - (float64(dist) / (float64(c.size) / 255))
					r, g, b, _ := c.col.RGBA()
					col = color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(alph)}
				} else if opacity != 100 {
					alph := opacity * 255 / 100
					r, g, b, _ := c.col.RGBA()
					col = color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(alph)}
				} else {
					col = c.col
				}

				if filled {
					if dist < c.size {
						circleimg.Set(x, y, col)
					}
				} else {
					if dist < c.size && dist > (c.size-bordersize) {
						circleimg.Set(x, y, col)
					}
				}
			}
		}

		dst := image.Rect(c.x-c.size, c.y-c.size, c.x+c.size, c.y+c.size)
		draw.Draw(img, dst, circleimg, image.Point{0, 0}, draw.Over)
	}

	return img
}

type Ray struct {
	col   color.Color
	x, y  int // Middle point
	angle int // 0-180
	size  int
}

// For sorting rays by size
type rayBySize []Ray

func (a rayBySize) Len() int           { return len(a) }
func (a rayBySize) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a rayBySize) Less(i, j int) bool { return a[i].size < a[j].size }

func Rays(colors []color.Color, w int, h int, size int, sizevar int, evendist bool, centered bool, large2small bool) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))

	rays := make([]Ray, 0)

	spacing := 360 / len(colors)
	current_angle := 0

	xpos := w / 2
	ypos := h / 2

	bg := colors[0]

	for _, c := range colors {
		// Do not create ray with background color
		if c == bg {
			continue
		}
		var ray Ray
		if !centered {
			xpos = rand.Intn(w)
			ypos = rand.Intn(h)
		}
		if !evendist {
			current_angle = rand.Intn(360)
		}
		ray = Ray{c, xpos, ypos, current_angle, randMinMax(size-sizevar, size+sizevar)}

		if evendist {
			current_angle += spacing + ray.size
		}
		rays = append(rays, ray)
	}

	if large2small {
		sort.Sort(sort.Reverse(rayBySize(rays)))
	}

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			img.Set(x, y, bg)
			for _, r := range rays {
				deltaX := float64(x - r.x)
				deltaY := float64(y - r.y)
				angle := math.Atan(deltaY/deltaX) * 180 / math.Pi
				if angle < 0 {
					angle += 360
				}
				if int(math.Abs(float64(int(angle)-r.angle))) < r.size {
					img.Set(x, y, r.col)
				}
			}
		}
	}
	return img
}

type Line struct {
	col      color.Color
	position int
	size     int
}

func Lines(colors []color.Color, w int, h int, size int, sizevar int, horizontal bool, equalspacing bool, spacingsize int, offset int) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, w, h))
	var maxsize int
	if horizontal {
		maxsize = h
	} else {
		maxsize = w
	}

	currentposition := offset
	spacing := spacingsize

	lines := make([]Line, 0)
	bg := colors[0]

	for _, c := range colors {
		// Do not create line with background color
		if c == bg {
			continue
		}
		line := Line{c, currentposition, randMinMax(size-sizevar, size+sizevar)}
		lines = append(lines, line)
		if !equalspacing {
			spacing = rand.Intn(maxsize / 16)
		}
		currentposition += line.size + spacing
	}

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			img.Set(x, y, bg)
			for _, l := range lines {
				var pixelpos int
				if horizontal {
					pixelpos = y
				} else {
					pixelpos = x
				}

				if pixelpos > l.position && pixelpos < l.position+l.size {
					img.Set(x, y, l.col)
				}
			}
		}
	}

	return img
}

func randomImage(colors []color.Color, w int, h int) image.Image {
	switch rand.Intn(3) {
	case 0:
		return Circles(colors, w, h, rand.Intn(w/2), rand.Intn(w/2), randBool(), randBool(), randBool(), rand.Intn(20), randBool(), 100)
	case 1:
		return Rays(colors, w, h, rand.Intn(h/32)+1, rand.Intn(h/32), randBool(), true, randBool())
	case 2:
		return Lines(colors, w, h, rand.Intn(h/32)+1, rand.Intn(h/32), randBool(), randBool(), rand.Intn(h/32), rand.Intn(h/2)+1)
	}
	return nil
}
