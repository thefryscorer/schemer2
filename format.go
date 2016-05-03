package main

import (
	"image/color"
)

type inputFunction (func(filename string) ([]color.Color, error))
type outputFunction (func([]color.Color) string)

type Format struct {
	friendlyName string
	flagName     string
	output       outputFunction
	input        inputFunction
}

var formats = []Format{
	{
		friendlyName: "Colors in Plain Text",
		flagName:     "colors",
		output:       printColors,
	},
	{
		friendlyName: "Image",
		flagName:     "img",
		input:        colorsFromImage,
	},
	{
		friendlyName: "XFCE4Terminal",
		flagName:     "xfce",
		input:        inputXfce,
		output:       printXfce,
	},
	{
		friendlyName: "LilyTerm",
		flagName:     "lilyterm",
		output:       printLilyTerm,
		input:        inputLilyTerm,
	},
	{
		friendlyName: "Termite",
		flagName:     "termite",
		input:        inputTermite,
		output:       printTermite,
	},
	{
		friendlyName: "Terminator",
		flagName:     "terminator",
		input:        inputTerminator,
		output:       printTerminator,
	},
	{
		friendlyName: "ROXTerm",
		flagName:     "roxterm",
		output:       printRoxTerm,
	},
	{
		friendlyName: "rxvt/xterm/aterm",
		flagName:     "xterm",
		input:        inputXterm,
		output:       printXterm,
	},
	{
		friendlyName: "Konsole",
		flagName:     "konsole",
		output:       printKonsole,
	},
	{
		friendlyName: "iTerm2",
		flagName:     "iterm2",
		output:       printITerm2,
	},
	{
		friendlyName: "urxvt",
		flagName:     "urxvt",
		input:        inputXterm,
		output:       printURxvt,
	},
	{
		friendlyName: "Chrome Shell",
		flagName:     "chrome",
		output:       printChrome,
	},
	{
		friendlyName: "OS X Terminal",
		flagName:     "osxterminal",
		output:       printOSXTerminal,
	},
	{
		friendlyName: "Gnome Terminal (dconf)",
		flagName:     "gnome-terminal",
		output:       printGnomeDConf,
	},
}
