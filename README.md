# golab

[![Actions Status](https://github.com/icza/golab/workflows/Go/badge.svg)](https://github.com/icza/golab/actions)

_This the reincarnation of my [gophergala/golab](https://github.com/gophergala/golab) game._

## Introduction

**Gopher's Labyrinth** (or **GoLab**) is a 2D Labyrinth game where you control Gopher
(who else) and your goal is to get to the Exit point.
But beware of the bloodthirsty Bulldogs, the ancient enemies of gophers who are endlessly roaming the Labyrinth!

Controlling Gopher is very easy: just click with your left mouse button to where you want to move
(there must be a free straight line to it). You may queue multiple target points forming a path.
Right click clears the path. You may also use the arrow keys on your keyboard.

You may try out the game in your browser if it supports WebAssembly and WebGL here: https://icza.github.io/golab/

![Screenshot](https://raw.githubusercontent.com/icza/golab/master/screenshot-golab.png)

## Under the hood

GoLab is written completely in [Go](https://golang.org). Go 1.13 or newer is required.
The user interface and input handling is done with the [gioui](https://gioui.org) library, utilized in the `view` package.

The game model and game logic is placed in the `engine` package.

## How to get it or install it

Of course in the "go way". You may quickly test it by initializing a new module in a folder by running:

    go mod init test
    
And then run GoLab with:

    go run github.com/icza/golab/cmd/golab

Or try it in your browser:  https://icza.github.io/golab/

## LICENSE

See [LICENSE](https://github.com/icza/golab/blob/master/LICENSE).

GoLab's Gopher is a derivative work based on the Go gopher which was designed by Renee French. (http://reneefrench.blogspot.com/).
Licensed under the Creative Commons 3.0 Attributions license.

The source of other images can be found in the [_images_/source.txt](https://github.com/icza/golab/blob/master/_images/source.txt) file.
