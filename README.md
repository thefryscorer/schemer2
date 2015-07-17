# UNDER CONSTRUCTION

Various aspects of this program are still being developed. Command line usage is quite probably going to change somewhat from as it is now. 

Terminal Colorscheme Generator
==============================

## Screenshot 
![Screenshot](http://i.imgur.com/63R70Iz.png)
[More Screenshots](http://imgur.com/a/Kz9Av)

## Installation 

### Short version

> go get github.com/thefryscorer/schemer2

### Long Version

#### Installing and configuring Go
To build this program, you will need to have Go installed and properly configured. After installing the Go package, you will need to configure a GOPATH. This is a directory in which Go will keep its programs and source files. I recommend making the GOPATH directory in your home folder. If your GOPATH is in your root directory a kitten will die. 

> mkdir ~/Go

You will also need to set the GOPATH variable so that Go knows where to put things. You can do this by running:

> export GOPATH=$HOME/Go

NOTE: You don't need to (and shouldn't) set the $GOROOT variable. This is handled for you and you shouldn't mess with it.

#### Installing schemer
You should now be able to install schemer using the command:

> go get github.com/thefryscorer/schemer2

And it will be built in your GOPATH directory, in a subdirectory named 'bin'. To run it, you can either add $GOPATH/bin to your system path and run it as you would any other command. Or cd into the bin directory and run it with:

> ./schemer2

## Usage 

#### Reading from terminal config and outputting to image
> schemer2 -in=xterm:$HOME/.Xresources -outputImage=out.png

#### Reading from that image and outputting terminal config (lilyterm)
> schemer2 -in=img:out.png -out=lilyterm

#### Reading from Xresources and outputting in termite format
> schemer2 -in=xterm:$HOME/.Xresources -out=termite

#### Getting colors from image, and outputting a new image
> schemer2 -in=img:old.png -outputImage=new.png


## Features 

- Outputs configuration in several different formats for different terminals.
- Configurable color difference threshold
- Configurable minimum and maximum brightness value

## Supported input formats

- Images (png, jpeg)
- Xfce4 terminalrc
- Lilyterm config
- Terminator config
- Termite config
- Xterm/URXvt and variants

## Supported output formats

- Colours in just plain text (default)
- Konsole
- xterm/rxvt/aterm
- urxvt
- iTerm2
- XFCE Terminal
- Roxterm
- LilyTerm
- Terminator
- Chrome Shell
- OS X Terminal
