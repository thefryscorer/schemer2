package main

import (
	"encoding/hex"
	"errors"
	"image/color"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

func readFile(filename string) (string, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	config := string(bytes[:])
	return config, nil
}

func parseColor(c string) (color.Color, error) {
	// Takes in a string of the format #FFFFFF or #FFFFFFFFFFFF and returns a color
	// Remove leading #
	c = strings.TrimPrefix(c, "#")
	// Convert hexadecimal string to array of bytes.
	bytes, err := hex.DecodeString(c)
	if err != nil {
		return nil, err
	}

	// Covert to array of uint8s
	uints := []uint8(bytes[:])

	// Take 0,1,2 indexes for a 6 character string, and 0,2,4 indexes for 12 characters
	if len(uints) == 6 {
		return color.NRGBA{uints[0], uints[2], uints[4], 255}, nil
	} else if len(uints) == 3 {
		return color.NRGBA{uints[0], uints[1], uints[2], 255}, nil
	}

	return nil, errors.New("Could not parse color: " + c)
}

func inputXfce(filename string) ([]color.Color, error) {
	// Read in file
	config, err := readFile(filename)
	if err != nil {
		return nil, err
	}

	// Split into lines
	lines := strings.Split(config, "\n")

	// Remove all spaces
	for i, l := range lines {
		lines[i] = strings.Replace(l, " ", "", -1)
	}

	// Find line containing color palette
	colorpalette := ""
	for _, l := range lines {
		if strings.HasPrefix(l, "ColorPalette") {
			colorpalette = l
		}
	}
	if colorpalette == "" {
		return nil, errors.New("ColorPalette not found in XFCE4 Terminal input")
	}

	// Get colors from palette
	colorpalette = strings.TrimPrefix(colorpalette, "ColorPalette=")
	// Trim trailing semicolon and spaces
	colorpalette = strings.TrimRight(colorpalette, "; ")
	// Split by semicolons
	colorStrings := strings.Split(colorpalette, ";")

	colors := make([]color.Color, 0)

	for _, c := range colorStrings {
		col, err := parseColor(c)
		if err != nil {
			return nil, err
		}
		colors = append(colors, col)
	}

	return colors, nil
}

func inputLilyTerm(filename string) ([]color.Color, error) {
	colors := make([]color.Color, 0)

	// Read in file
	config, err := readFile(filename)
	if err != nil {
		return nil, err
	}

	// Split into lines
	lines := strings.Split(config, "\n")

	// Remove all spaces
	for i, l := range lines {
		lines[i] = strings.Replace(l, " ", "", -1)
	}

	// For all 16 colors (Color1, Color2...), search for each.
	// TODO: Support for 16 bit, and "black"
	for i := 0; i < 16; i++ {
		for _, l := range lines {
			prefix := "Color"
			prefix += strconv.Itoa(i)
			prefix += "="
			if strings.HasPrefix(l, prefix) {
				// Trim Prefix
				hexstring := strings.TrimPrefix(l, prefix)

				col, err := parseColor(hexstring)
				if err != nil {
					return nil, err
				}

				colors = append(colors, col)

			}
		}
	}

	return colors, nil
}

func inputTermite(filename string) ([]color.Color, error) {
	colors := make([]color.Color, 0)

	// Read in file
	config, err := readFile(filename)
	if err != nil {
		return nil, err
	}

	// Split into lines
	lines := strings.Split(config, "\n")

	// Remove all spaces
	for i, l := range lines {
		lines[i] = strings.Replace(l, " ", "", -1)
	}

	// For all 16 colors (Color1, Color2...), search for each.
	for i := 0; i < 16; i++ {
		for _, l := range lines {
			prefix := "color"
			prefix += strconv.Itoa(i)
			prefix += "="
			if strings.HasPrefix(l, prefix) {
				// Trim Prefix
				hexstring := strings.TrimPrefix(l, prefix)

				col, err := parseColor(hexstring)
				if err != nil {
					return nil, err
				}

				colors = append(colors, col)
			}
		}
	}

	return colors, nil
}

func inputTerminator(filename string) ([]color.Color, error) {
	// Read in file
	config, err := readFile(filename)
	if err != nil {
		return nil, err
	}

	// Split into lines
	lines := strings.Split(config, "\n")

	// Remove all spaces
	for i, l := range lines {
		lines[i] = strings.Replace(l, " ", "", -1)
	}

	// Find line containing color palette
	colorpalette := ""
	for _, l := range lines {
		if strings.HasPrefix(l, "palette") {
			colorpalette = l
		}
	}
	if colorpalette == "" {
		return nil, errors.New("ColorPalette not found in XFCE4 Terminal input")
	}

	// Get colors from palette
	colorpalette = strings.TrimPrefix(colorpalette, "palette=\"")
	colorpalette = strings.TrimSuffix(colorpalette, "\"")

	colorStrings := strings.Split(colorpalette, ":")

	colors := make([]color.Color, 0)

	for _, c := range colorStrings {
		col, err := parseColor(c)
		if err != nil {
			return nil, err
		}
		colors = append(colors, col)
	}
	return colors, nil

}

func inputXterm(filename string) ([]color.Color, error) {
	// Read in file
	config, err := readFile(filename)
	if err != nil {
		return nil, err
	}

	// Split into lines
	lines := strings.Split(config, "\n")

	// Remove all spaces
	for i, l := range lines {
		lines[i] = strings.Replace(l, " ", "", -1)
	}

	colorlines := make([]string, 0)
	// Search for lines containing color information
	re := regexp.MustCompile("[\\*]?[URXvurxterm]*[\\*.]+color[0-9]*")
	for _, l := range lines {
		if len(re.FindAllString(l, 1)) != 0 {
			colorlines = append(colorlines, l)
		}
	}

	// Extract and parse colors
	// TODO: Sort by number first?
	colors := make([]color.Color, 0)
	for _, l := range colorlines {
		// Assuming the color to be the rightmost half of the last colon
		splits := strings.Split(l, ":")
		colorstring := splits[len(splits)-1]
		col, err := parseColor(colorstring)
		if err != nil {
			return nil, err
		}
		colors = append(colors, col)
		if len(colorlines) < 16 && i < 16 - len(colorlines) {
			colors = append(colors, col)
		}
	}

	return colors, nil
}
