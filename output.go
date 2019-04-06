package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"image/color"
	"strconv"
)

func printXfce(colors []color.Color) string {
	output := ""
	output += "ColorPalette="
	for _, c := range colors {
		bytes := []byte{byte(c.(color.NRGBA).R), byte(c.(color.NRGBA).R), byte(c.(color.NRGBA).G), byte(c.(color.NRGBA).G), byte(c.(color.NRGBA).B), byte(c.(color.NRGBA).B)}
		output += "#"
		output += hex.EncodeToString(bytes)
		output += ";"
	}
	output += "\n"

	return output
}

func printLilyTerm(colors []color.Color) string {
	output := ""
	for i, c := range colors {
		cc := c.(color.NRGBA)
		bytes := []byte{byte(cc.R), byte(cc.G), byte(cc.B)}
		output += "Color"
		output += strconv.Itoa(i)
		output += " = "
		output += "#"
		output += hex.EncodeToString(bytes)
		output += "\n"
	}
	return output
}

func printTermite(colors []color.Color) string {
	output := ""
	for i, c := range colors {
		cc := c.(color.NRGBA)
		bytes := []byte{byte(cc.R), byte(cc.G), byte(cc.B)}
		output += "color"
		output += strconv.Itoa(i)
		output += " = "
		output += "#"
		output += hex.EncodeToString(bytes)
		output += "\n"
	}
	return output
}

func printTerminator(colors []color.Color) string {
	output := "palette = \""
	for i, c := range colors {
		cc := c.(color.NRGBA)
		bytes := []byte{byte(cc.R), byte(cc.G), byte(cc.B)}
		if i < len(colors)-1 {
			output += "#"
			output += hex.EncodeToString(bytes)
			output += ":"
		} else if i == len(colors)-1 {
			output += "#"
			output += hex.EncodeToString(bytes)
			output += "\"\n"
		}
	}
	return output
}

func printXterm(colors []color.Color) string {
	output := ""
	output += "! Terminal colors"
	output += "\n"
	for i, c := range colors {
		cc := c.(color.NRGBA)
		bytes := []byte{byte(cc.R), byte(cc.G), byte(cc.B)}
		output += "*color"
		output += strconv.Itoa(i)
		output += ": #"
		output += hex.EncodeToString(bytes)
		output += "\n"
	}

	return output
}

func printKonsole(colors []color.Color) string {
	output := ""
	for i, c := range colors {
		cc := c.(color.NRGBA)
		output += "[Color"
		if i > 7 {
			output += strconv.Itoa(i - 8)
			output += "Intense"
		} else {
			output += strconv.Itoa(i)
		}
		output += "]\n"
		output += "Color="
		output += strconv.Itoa(int(cc.R)) + ","
		output += strconv.Itoa(int(cc.G)) + ","
		output += strconv.Itoa(int(cc.B)) + "\n"
		output += "Transparency=false\n\n"
	}

	return output
}

func printRoxTerm(colors []color.Color) string {
	output := "[roxterm colour scheme]\n"
	output += "pallete_size=16\n"

	for i, c := range colors {
		cc := c.(color.NRGBA)
		bytes := []byte{byte(cc.R), byte(cc.G), byte(cc.B)}
		output += "color"
		output += strconv.Itoa(i)
		output += " = "
		output += "#"
		output += hex.EncodeToString(bytes)
		output += "\n"
	}

	return output
}

func printITerm2(colors []color.Color) string {
	output := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"
	output += "<!DOCTYPE plist PUBLIC \"-//Apple//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\">\n"
	output += "<plist version=\"1.0\">\n"
	output += "<dict>\n"
	for i, c := range colors {
		cc := c.(color.NRGBA)
		output += "\t<key>Ansi "
		output += strconv.Itoa(i)
		output += " Color</key>\n"
		output += "\t<dict>\n"
		output += "\t\t<key>Blue Component</key>\n"
		output += "\t\t<real>"
		output += strconv.FormatFloat(float64(cc.B)/255, 'f', 17, 64)
		output += "</real>\n"
		output += "\t\t<key>Green Component</key>\n"
		output += "\t\t<real>"
		output += strconv.FormatFloat(float64(cc.G)/255, 'f', 17, 64)
		output += "</real>\n"
		output += "\t\t<key>Red Component</key>\n"
		output += "\t\t<real>"
		output += strconv.FormatFloat(float64(cc.R)/255, 'f', 17, 64)
		output += "</real>\n"
		output += "\t</dict>\n"
	}
	output += "</dict>\n"
	output += "</plist>\n"
	return output
}

func printURxvt(colors []color.Color) string {
	output := ""
	for i, c := range colors {
		cc := c.(color.NRGBA)
		bytes := []byte{byte(cc.R), byte(cc.G), byte(cc.B)}
		output += "URxvt*color"
		output += strconv.Itoa(i)
		output += ": "
		output += "#"
		output += hex.EncodeToString(bytes)
		output += "\n"
	}
	return output
}

func printColors(colors []color.Color) string {
	output := ""
	for _, c := range colors {
		cc := c.(color.NRGBA)
		bytes := []byte{byte(cc.R), byte(cc.G), byte(cc.B)}
		output += "#"
		output += hex.EncodeToString(bytes)
		output += "\n"
	}
	return output
}
func printChrome(colors []color.Color) string {
	output := "{"
	for i, c := range colors {
		cc := c.(color.NRGBA)
		bytes := []byte{byte(cc.R), byte(cc.G), byte(cc.B)}
		output += " \""
		output += strconv.Itoa(i)
		output += "\""
		output += ": "
		output += " \""
		output += "#"
		output += hex.EncodeToString(bytes)
		output += "\" "
		if i != len(colors)-1 {
			output += ", "
		}
	}
	output += "}\n"
	return output
}

func printOSXTerminal(colors []color.Color) string {
	// The plist that is used by OS X's Terminal to store colours. Normally,
	// Terminal stores the colours in a base64 encoded binary plist but it'll
	// happily read base64 encoded xml plists which makes things easier.
	const OSXSerializedNSColorTemplate = `<?xml version="1.0" encoding="UTF-8"?><!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd"><plist version="1.0"><dict><key>$archiver</key><string>NSKeyedArchiver</string><key>$objects</key><array><string>$null</string><dict><key>$class</key><dict><key>CF$UID</key><integer>2</integer></dict><key>NSColorSpace</key><integer>1</integer><key>NSRGB</key><data>%s</data></dict><dict><key>$classes</key><array><string>NSColor</string><string>NSObject</string></array><key>$classname</key><string>NSColor</string></dict></array><key>$top</key><dict><key>root</key><dict><key>CF$UID</key><integer>1</integer></dict></dict><key>$version</key><integer>100000</integer></dict></plist>`
	OSXColorNames := map[int]string{
		0: "Black",
		1: "Red",
		2: "Green",
		3: "Yellow",
		4: "Blue",
		5: "Magenta",
		6: "Cyan",
		7: "White",
	}

	output := "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n"
	output += "<!DOCTYPE plist PUBLIC \"-//Apple//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\">\n"
	output += "<plist version=\"1.0\">\n"
	output += "<dict>\n"
	for i, c := range colors {
		cc := c.(color.NRGBA)
		output += "\t<key>ANSI"
		if i > 7 {
			output += "Bright" + OSXColorNames[i-8]
		} else {
			output += OSXColorNames[i]
		}
		output += "Color</key>\n"
		output += "\t<data>\n"
		rgbColorString := fmt.Sprintf("%.10f %.10f %.10f", float64(cc.R)/255, float64(cc.G)/255, float64(cc.B)/255)
		serializedColor := fmt.Sprintf(OSXSerializedNSColorTemplate, base64.StdEncoding.EncodeToString([]byte(rgbColorString)))
		output += "\t" + base64.StdEncoding.EncodeToString([]byte(serializedColor))
		output += "\n\t</data>\n"
	}

	output += "\t<key>type</key>\n" // Need this key or Terminal says the file is corrupt
	output += "\t<string>Window Settings</string>\n"
	output += "</dict>\n"
	output += "</plist>\n"
	return output
}

func printGnomeDConf(colors []color.Color) string {
	output := "#!/usr/bin/env bash\npalette=\"["
	for i, c := range colors {
		cc := c.(color.NRGBA)
		bytes := []byte{byte(cc.R), byte(cc.G), byte(cc.B)}
		output += "'"
		output += "#"
		output += hex.EncodeToString(bytes)
		output += "'"
		if i < len(colors)-1 {
			output += ","
		}
	}
	output += "]\""
	output += "\n"
	output += "default=$(dconf read /org/gnome/terminal/legacy/profiles:/default | sed -e \"s/'//g\")"
	output += "\n"
	output += "dconf write /org/gnome/terminal/legacy/profiles:/:$default/palette \"$palette\""
	output += "\n"
	return output
}

func printKittyTerm(colors []color.Color) string {
	output := ""
	for i, c := range colors {
		cc := c.(color.NRGBA)
		bytes := []byte{byte(cc.R), byte(cc.G), byte(cc.B)}
		output += "color"
		output += strconv.Itoa(i)
		output += "\t#"
		output += hex.EncodeToString(bytes)
		output += "\n"
	}

	return output
}
