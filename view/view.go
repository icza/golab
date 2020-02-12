// Package view is the view of the game. It handles user input and presents the game to the user.
package view

import (
	"image"
	"image/color"
	"image/draw"
	"log"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/icza/golab/ctrl"
	"github.com/icza/golab/model"
)

const (
	controlsHeight = 50
	viewWidth      = 700
	viewHeight     = 700
	// WindowWidth is the suggested window width
	WindowWidth = viewWidth
	// WindowHeight is the suggested window height
	WindowHeight = controlsHeight + viewHeight
)

func init() {
	gofont.Register()
}

// View displays the game window and handles user input.
type View struct {
	engine *ctrl.Engine
	w      *app.Window
	th     *material.Theme
	gtx    *layout.Context

	// New Game button model
	newGameBtn *widget.Button

	// Difficulty options
	diffOpt *options
	// Lab size options
	labSizeOpt *options
	// Speed options
	speedOpt *options

	// "static" imageOps
	imgOpGophers  []paint.ImageOp
	imgOpDead     paint.ImageOp
	imgOpBulldogs []paint.ImageOp
	imgOpMarker   paint.ImageOp
	imgOpExit     paint.ImageOp
	imgOpWon      paint.ImageOp

	// gameCounter for the cached data
	gameCounter int
	// cached ImageOp of the whole labyrinth (only the blocks)
	labImgOp paint.ImageOp
}

// New returns a new View.
func New(engine *ctrl.Engine, w *app.Window) *View {
	v := &View{
		engine:      engine,
		w:           w,
		th:          material.NewTheme(),
		gtx:         layout.NewContext((w.Queue())),
		newGameBtn:  new(widget.Button),
		imgOpDead:   paint.NewImageOp(imgDead),
		imgOpMarker: paint.NewImageOp(imgMarker),
		imgOpExit:   paint.NewImageOp(imgExit),
		imgOpWon:    paint.NewImageOp(imgWon),
	}

	for _, img := range imgGophers {
		v.imgOpGophers = append(v.imgOpGophers, paint.NewImageOp(img))
	}
	for _, img := range imgBulldogs {
		v.imgOpBulldogs = append(v.imgOpBulldogs, paint.NewImageOp(img))
	}

	v.diffOpt = newOptions(v, "Difficulty", ctrl.Difficulties, ctrl.DifficultyDefaultIdx)
	v.labSizeOpt = newOptions(v, "Lab size", ctrl.LabSizes, ctrl.LabSizeDefaultIdx)
	v.speedOpt = newOptions(v, "Speed", ctrl.Speeds, ctrl.SpeedDefaultIdx)

	return v
}

// Loop starts handing user input and frame redraws.
// This function returns only if the user closes the app.
func (v *View) Loop() {
	for e := range v.w.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			v.drawFrame(e)
		case system.DestroyEvent:
			log.Println("Goodbye!")
		}
	}
}

// drawFrame draws a frame of the window.
func (v *View) drawFrame(e system.FrameEvent) {
	gtx := v.gtx

	gtx.Reset(e.Config, e.Size)

	// Handle button clicks
	for v.newGameBtn.Clicked(v.gtx) {
		v.engine.NewGame(ctrl.GameConfig{
			Difficulty: v.diffOpt.selected().(*ctrl.Difficulty),
			LabSize:    v.labSizeOpt.selected().(*ctrl.LabSize),
			Speed:      v.speedOpt.selected().(*ctrl.Speed),
		})
	}
	v.diffOpt.handleInput()
	v.labSizeOpt.handleInput()
	v.speedOpt.handleInput()

	v.drawControls()
	v.drawLab()

	e.Frame(gtx.Ops)
}

// drawControls draws the control and setup widgets.
func (v *View) drawControls() {
	th, gtx := v.th, v.gtx

	layout.N.Layout(gtx, func() {
		layout.UniformInset(unit.Px(5)).Layout(gtx, func() {
			layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func() {
					layout.Inset{Left: unit.Px(10), Right: unit.Px(10)}.Layout(gtx, func() {
						b := th.Button("New Game")
						b.Background = color.RGBA{R: 20, G: 130, B: 20, A: 255}
						b.Layout(gtx, v.newGameBtn)
					})
				}),
				layout.Rigid(v.diffOpt.layout),
				layout.Rigid(v.labSizeOpt.layout),
				layout.Rigid(v.speedOpt.layout),
			)
		})
	})
}

// drawLab draws the labyrinth.
func (v *View) drawLab() {
	m := v.engine.Model
	m.RLock()
	defer m.RUnlock()

	gtx := v.gtx

	var stack op.StackOp
	stack.Push(gtx.Ops)
	defer stack.Pop()

	// Center lab:
	displayWidth, displayHeight := viewWidth, viewHeight
	if labWidth := m.Cols * ctrl.BlockSize; labWidth < displayWidth {
		displayWidth = labWidth
	}
	if labHeight := m.Rows * ctrl.BlockSize; labHeight < displayHeight {
		displayHeight = labHeight
	}
	op.TransformOp{}.Offset(f32.Point{
		X: float32((gtx.Constraints.Width.Max - displayWidth) / 2),
		Y: controlsHeight,
	}).Add(gtx.Ops)
	clip.Rect{Rect: f32.Rectangle{Max: f32.Point{
		X: float32(displayWidth),
		Y: float32(displayHeight),
	}}}.Op(gtx.Ops).Add(gtx.Ops)

	// First the blocks:
	v.ensureLabImgOp()
	v.drawImg(v.labImgOp, 0, 0)

	// Now objects in the lab:
	// Gopher:
	if m.Dead {
		v.drawObj(v.imgOpDead, m.Gopher)
	} else {
		v.drawObj(v.imgOpGophers[m.Gopher.Dir], m.Gopher)
	}
	// Bulldogs:
	for _, bd := range m.Bulldogs {
		v.drawObj(v.imgOpBulldogs[bd.Dir], bd)
	}

	// TODO
}

// drawObj draws the given image of the given moving obj.
func (v *View) drawObj(iop paint.ImageOp, obj *model.MovingObj) {
	v.drawImg(iop, float32(obj.Pos.X-ctrl.BlockSize/2), float32(obj.Pos.Y-ctrl.BlockSize/2))
}

// drawImg draws the given image to the given position.
func (v *View) drawImg(iop paint.ImageOp, x, y float32) {
	var stack op.StackOp
	stack.Push(v.gtx.Ops)

	op.TransformOp{}.Offset(f32.Point{X: x, Y: y}).Add(v.gtx.Ops)

	img := v.th.Image(iop)
	img.Scale = 1
	img.Layout(v.gtx)

	stack.Pop()
}

// ensureLabImgOp ensures labImgOp is up-to-date
func (v *View) ensureLabImgOp() {
	m := v.engine.Model
	if v.gameCounter == m.Counter {
		// We have the lab image of the current game
		return
	}

	labImg := image.NewRGBA(image.Rect(0, 0, m.Cols*ctrl.BlockSize, m.Rows*ctrl.BlockSize))
	var r image.Rectangle
	for row := range m.Lab {
		r.Min.Y = row * ctrl.BlockSize
		r.Max.Y = r.Min.Y + ctrl.BlockSize
		for col, block := range m.Lab[row] {
			r.Min.X = col * ctrl.BlockSize
			r.Max.X = r.Min.X + ctrl.BlockSize
			src := imgBlocks[block]
			draw.Draw(labImg, r, src, image.Point{}, draw.Over)
		}
	}

	// Exit sign:
	r.Min = m.ExitPos
	r.Min = r.Min.Add(image.Point{-ctrl.BlockSize / 2, -ctrl.BlockSize / 2})
	r.Max = r.Min.Add(image.Point{ctrl.BlockSize, ctrl.BlockSize})
	draw.Draw(labImg, r, imgExit, image.Point{}, draw.Over)

	v.labImgOp = paint.NewImageOp(labImg)

	v.gameCounter = m.Counter
}
