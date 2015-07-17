package main

import (
	"flag"
	"fmt"
	"image/color"
	"image/png"
	"log"
	"os"
	"strings"
)

var (
	threshold     *int
	output        *string
	input         *string
	minBrightness *int
	maxBrightness *int
	imageout      *string
	imageWidth    *int
	imageHeight   *int
)

func usage() {
	fmt.Println("Usage: schemer2 [FLAGS] -in=[FORMAT]:[FILENAME] (-out=[FORMAT] | -outputImage=[FILENAME])")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	inSupport := "Format and filename of input file (eg \"xfce:~/.config/xfce4/terminal/terminalrc\"). Currently supported: \n"
	outSupport := "Format to output colors as. Currently supported: \n"
	for _, f := range formats {
		if f.output != nil {
			outSupport += strings.Join([]string{"    ", f.friendlyName, ":", f.flagName, "\n"}, " ")
		}

		if f.input != nil {
			inSupport += strings.Join([]string{"    ", f.friendlyName, ":", f.flagName, "\n"}, " ")
		}
	}

	threshold = flag.Int("t", 50, "Threshold for minimum color difference (image input only)")
	output = flag.String("out", "", outSupport)
	input = flag.String("in", "", inSupport)
	minBrightness = flag.Int("minBright", 0, "Minimum brightness for colors (image input only)")
	maxBrightness = flag.Int("maxBright", 200, "Maximum brightness for colors (image input only)")
	imageout = flag.String("outputImage", "", "Create image from colors, and save to this file")
	imageHeight = flag.Int("h", 1080, "Height of output image")
	imageWidth = flag.Int("w", 1920, "Width of output image")

	flag.Usage = usage
	flag.Parse()
	if *input == "" {
		usage()
		os.Exit(2)
	}
	if *minBrightness > 255 || *maxBrightness > 255 {
		fmt.Print("Minimum and maximum brightness must be an integer between 0 and 255.\n")
		os.Exit(2)
	}
	if *threshold > 255 {
		fmt.Print("Threshold should be an integer between 0 and 255.\n")
		os.Exit(2)
	}

	if *imageWidth < 100 || *imageHeight < 100 {
		log.Fatal("Minimum resolution of image output is 100x100")
	}

	// Determine format and filename
	// And get colors from file using specified format
	format := strings.SplitN(*input, ":", 2)[0]
	filename := strings.SplitN(*input, ":", 2)[1]

	formatInMatch := false
	var colors []color.Color
	var err error
	for _, f := range formats {
		if format == f.flagName {
			if f.input == nil {
				fmt.Printf("Unrecognised input format: %v \n", format)
				return
			}
			colors, err = f.input(filename)
			if err != nil {
				fmt.Print(err, "\n")
				return
			}
			formatInMatch = true
			break
		}
	}
	if !formatInMatch {
		fmt.Printf("Did not recognise format %v. \n", *input)
		return
	}

	// Ensure there are 16 colors
	if len(colors) > 16 {
		colors = colors[:16]
	} else if len(colors) < 16 {
		// TODO: Should this just be a warning (for cases where only 8 colors are defined?)
		log.Fatal("Less than 16 colors. Aborting.")
	}

	// Output the configuration specified
	if !(*output == "") {
		formatOutMatch := false
		for _, f := range formats {
			if *output == f.flagName {
				if f.output == nil {
					fmt.Printf("Unrecognised output format: %v \n", format)
					return
				}
				fmt.Print(f.output(colors))
				formatOutMatch = true
				break
			}
		}
		if !formatOutMatch {
			fmt.Printf("Did not recognise format %v. \n", *output)
		}
	}

	if *imageout != "" {
		file, err := os.OpenFile(*imageout, os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		img := imageFromColors(colors, *imageWidth, *imageHeight) // TODO

		png.Encode(file, img)
	}
}
