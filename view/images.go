package view

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"

	"github.com/icza/golab/ctrl"
	"github.com/icza/golab/model"
)

//go:generate go run _generate-embedded-imgs/main.go

// Tells if the embedded images are to be used. If false, images from files will be loaded.
const useEmbeddedImages = true

// imgGopher holds images of Gopher for each direction, each has zero Min point
var imgGophers = make([]*image.RGBA, model.DirCount)

// imgDead is the Dead Gopher image.
var imgDead *image.RGBA

// imgBulldog holds images of a Bulldog for each direction, each has zero Min point
var imgBulldogs = make([]*image.RGBA, model.DirCount)

// imgBlocks holds images of labyrinth blocks for each type, each has zero Min point
var imgBlocks = make([]image.Image, model.BlockCount)

// imgMarker is the image of the path marker
var imgMarker *image.RGBA

// imgExit is the image of the exit sign
var imgExit *image.RGBA

// imgWon is the image of the winning sign
var imgWon *image.RGBA

func init() {
	for dir := model.Dir(0); dir < model.DirCount; dir++ {
		// Load Gopher images
		imgGophers[dir] = loadImg(fmt.Sprintf("gopher-%s.png", dir), true)
		// Load Bulldog images
		imgBulldogs[dir] = loadImg(fmt.Sprintf("bulldog-%s.png", dir), true)
	}

	imgBlocks[model.BlockEmpty] = image.NewUniform(color.RGBA{A: 0xff})
	imgBlocks[model.BlockWall] = loadImg("wall.png", true)
	imgDead = loadImg("gopher-dead.png", true)
	imgExit = loadImg("door.png", true)

	imgMarker = loadImg("marker.png", false)
	imgWon = loadImg("won.png", false)
}

// loadImg loads a PNG image from the specified file, and converts it to image.RGBA and makes sure image has zero Min point.
// This function only used during development as the result contains the images embedded.
// blockSize tells if the image must be of the size of a block (else panics).
func loadImg(name string, blockSize bool) *image.RGBA {
	var data []byte
	var err error

	if useEmbeddedImages {
		data, err = base64.StdEncoding.DecodeString(base64Imgs[name])
	} else {
		data, err = ioutil.ReadFile(name)
	}
	if err != nil {
		panic(err)
	}
	return decodeImg(data, blockSize)
}

// decodeImg decodes an image from the specified data which must be of PNG format.
// blockSize tells if the image must be of the size of a block (else panics).
func decodeImg(data []byte, blockSize bool) *image.RGBA {
	src, err := png.Decode(bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}

	// Convert to image.RGBA, also make sure result image has zero Min point
	b := src.Bounds()
	if blockSize && (b.Dx() != ctrl.BlockSize || b.Dy() != ctrl.BlockSize) {
		panic("Invalid image size!")
	}

	img := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(img, src.Bounds(), src, b.Min, draw.Src)

	return img
}
