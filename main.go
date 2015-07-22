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

const (
	format_separator = "::"
)

var (
	outfile       *string
	infile        *string
	format_string *string

	// Image input options
	threshold     *int
	minBrightness *int
	maxBrightness *int

	// Generic image output options
	imageWidth   *int
	imageHeight  *int
	imageOutType *string // Eg, "random", "circles", "stripes", etc...
	imageOverlay *string

	// Circles image output options
	circlesSize                  *int
	circlesSizeVariance          *int
	circlesOverlap               *bool
	circlesDrawLargestToSmallest *bool
	circlesFilled                *bool
	circlesBorderSize            *int
	circlesBlur                  *bool
	circlesOpacity               *int

	// Ray image output options
	raysSize                  *int
	raysSizeVariance          *int
	raysDistributeEvenly      *bool
	raysCentered              *bool
	raysDrawLargestToSmallest *bool

	// Stripes image output options
	stripesSize         *int
	stripesSizeVariance *int
	stripesHorizontal   *bool
	stripesEqualSize    *bool
	stripesEvenSpacing  *bool
	stripesSpacing      *int
	stripesOffset       *int

	// Show advanced help
	advancedoptions *bool
)

func usage() {
	fmt.Println("Usage: schemer2 [FLAGS] -format [INPUTFORMAT]" + format_separator + "[OUTPUTFORMAT] -in [INPUTFILE] -out [OUTPUTFILE]")
}

func flags_usage() {
	usage()
	fmt.Printf("\n\n") // Spacing
	inputs_outputs()
	fmt.Printf("\n\n") // Spacing

	if !*advancedoptions {
		fmt.Println("Run with -help-advanced flag to show advanced options")
	} else {
		flag.PrintDefaults()
	}
}

func inputs_outputs() {
	inSupport := "Input formats:\n"
	outSupport := "Output formats:\n"
	for _, f := range formats {
		if f.output != nil {
			outSupport += strings.Join([]string{"    ", f.friendlyName, ":", f.flagName, "\n"}, " ")
		}

		if f.input != nil {
			inSupport += strings.Join([]string{"    ", f.friendlyName, ":", f.flagName, "\n"}, " ")
		}
	}
	// Special case for img output
	outSupport += "     Image output : img\n"

	fmt.Print(inSupport, "\n", outSupport)
}

func main() {
	infile = flag.String("in", "", "Input file")
	outfile = flag.String("out", "", "File to write output to.")
	format_string = flag.String("format", "", "Format of input and output. Eg. 'image"+format_separator+"xterm'")

	threshold = flag.Int("threshold", 50, "Threshold for minimum color difference (image input only)")
	minBrightness = flag.Int("minBright", 0, "Minimum brightness for colors (image input only)")
	maxBrightness = flag.Int("maxBright", 200, "Maximum brightness for colors (image input only)")

	imageHeight = flag.Int("height", 1080, "Height of output image")
	imageWidth = flag.Int("width", 1920, "Width of output image")
	imageOutTypeDesc := "Type of image to generate. Available options: \n"
	for _, t := range imageOutTypes {
		imageOutTypeDesc += "    "
		imageOutTypeDesc += t
		imageOutTypeDesc += "\n"
	}
	imageOutType = flag.String("imageOutType", "random", imageOutTypeDesc)
	imageOverlay = flag.String("imageOverlay", "", "Filename of image to draw on top of generated image (OS/Distro logo, etc...)")

	// Circles image output options
	circlesSize = flag.Int("circlesSize", 100, "Size of circles in output image")
	circlesSizeVariance = flag.Int("circlesSizeVariance", 50, "Maximum variance in circle size")
	circlesOverlap = flag.Bool("circlesOverlap", true, "Allow circles to overlap !!! Unimplemented !!!")
	circlesDrawLargestToSmallest = flag.Bool("circlesLargeToSmall", true, "Order circles z-index by size (smaller circles are drawn in front of larger circles)")
	circlesFilled = flag.Bool("circlesFilled", false, "Fill circles")
	circlesBlur = flag.Bool("circlesBlurred", false, "Blur circles")
	circlesOpacity = flag.Int("circlesOpacity", 100, "Opacity of circles")
	circlesBorderSize = flag.Int("circlesBorderSize", 10, "Border of circles when unfilled")

	// Ray image output options
	raysSize = flag.Int("raysSize", 16, "Size of rays in output image")
	raysSizeVariance = flag.Int("raysSizeVariance", 8, "Maximum variance in rays size")
	raysDistributeEvenly = flag.Bool("raysDistributeEvenly", false, "Distribute rays evenly")
	raysCentered = flag.Bool("raysCentered", true, "Center rays in middle")
	raysDrawLargestToSmallest = flag.Bool("raysLargeToSmall", false, "Order rays z-index by size (smaller rays are drawn on top of larger rays)")

	// Stripes image output options
	stripesSize = flag.Int("stripesSize", 6, "Size of stripes in output image")
	stripesSizeVariance = flag.Int("stripesSizeVariance", 3, "Maximum variance in stripes size")
	stripesHorizontal = flag.Bool("stripesHorizontal", false, "Draw stripes horizontally instead of vertically")
	stripesEvenSpacing = flag.Bool("stripesEvenSpacing", true, "Space all stripes evenly")
	stripesSpacing = flag.Int("stripesSpacing", 0, "Space stripes by this amount when spacing evenly")
	stripesOffset = flag.Int("stripesOffset", 0, "Offset stripes by this amount")

	advancedoptions = flag.Bool("help-advanced", false, "Show advanced command line options")

	flag.Usage = flags_usage
	flag.Parse()
	if *advancedoptions {
		flags_usage()
		os.Exit(1)
	}
	if *format_string == "" {
		fmt.Println("Input and output format must be specified using '-format' flag.")
		flags_usage()
		os.Exit(2)
	}
	if *infile == "" {
		fmt.Println("Input file must be provided using '-in' flag.")
		flags_usage()
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

	// Determine format and filenames
	// And get colors from file using specified format
	if len(strings.SplitN(*format_string, format_separator, 2)) < 2 {
		fmt.Println("Invalid format string. Separate input and output formats with: '" + format_separator + "'")
		flags_usage()
		os.Exit(2)
	}
	input_format := strings.SplitN(*format_string, format_separator, 2)[0]
	output_format := strings.SplitN(*format_string, format_separator, 2)[1]

	formatInMatch := false
	var colors []color.Color
	var err error
	for _, f := range formats {
		if input_format == f.flagName {
			if f.input == nil {
				fmt.Printf("Unrecognised input format: %v \n", input_format)
				return
			}
			colors, err = f.input(*infile)
			if err != nil {
				fmt.Print(err, "\n")
				return
			}
			formatInMatch = true
			break
		}
	}
	if !formatInMatch {
		fmt.Printf("Did not recognise format %v. \n", input_format)
		return
	}

	// Ensure there are 16 colors
	if len(colors) > 16 {
		colors = colors[:16]
	} else if len(colors) < 16 {
		// TODO: Should this just be a warning (for cases where only 8 colors are defined?)
		log.Fatal("Less than 16 colors. Aborting.")
	}

	// Output the configuration for terminal, or image
	if output_format != "img" {
		formatOutMatch := false
		for _, f := range formats {
			if output_format == f.flagName {
				if f.output == nil {
					fmt.Printf("Unrecognised output format: %v \n", output_format)
					return
				}
				result := f.output(colors)
				// If outfile is specified, write output to file
				// Otherwise, write to stdout.
				// TODO: Make it abundantly clear that the output is *only* the colors
				// and attempting to write directly to a config file will overwrite all other
				// data in the config file.
				if *outfile != "" {
					file, err := os.OpenFile(*outfile, os.O_CREATE|os.O_WRONLY, 0666)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					defer file.Close()
					err = file.Truncate(0)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					fmt.Fprint(file, result)
				} else {
					fmt.Printf(result)
				}
				formatOutMatch = true
				break
			}
		}
		if !formatOutMatch {
			fmt.Printf("Did not recognise format %v. \n", output_format)
		}
	} else {
		var file *os.File
		var err error
		if *outfile == "" {
			fmt.Println("Warning: Image output requested, yet no output file provided.")
			fmt.Println("Writing image data to /tmp/schemer_out.png")
			file, err = os.OpenFile("/tmp/schemer_out.png", os.O_CREATE|os.O_WRONLY, 0666)
		} else {
			file, err = os.OpenFile(*outfile, os.O_CREATE|os.O_WRONLY, 0666)
		}
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		img, err := imageFromColors(colors, *imageWidth, *imageHeight) // TODO
		if err != nil {
			fmt.Printf(err.Error())
			os.Exit(1)
		}

		png.Encode(file, img)
	}
}
