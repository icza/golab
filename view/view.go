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

	// gameCounter for the cached data
	gameCounter int
	// cached ImageOp of the whole labyrinth (only the blocks)
	labImgOp paint.ImageOp
}

// New returns a new View.
func New(engine *ctrl.Engine, w *app.Window) *View {
	v := &View{
		engine:     engine,
		w:          w,
		th:         material.NewTheme(),
		gtx:        layout.NewContext((w.Queue())),
		newGameBtn: new(widget.Button),
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

	th, gtx := v.th, v.gtx

	var stack op.StackOp
	stack.Push(gtx.Ops)
	defer stack.Pop()

	clip.Rect{Rect: f32.Rectangle{Max: f32.Point{X: 700, Y: 700}}}.Op(gtx.Ops).Add(gtx.Ops)

	// TODO

	// First the blocks:

	v.ensureLabImgOp()
	layout.Center.Layout(gtx, func() {
		img := th.Image(v.labImgOp)
		img.Scale = 1
		img.Layout(gtx)
	})

	paint.ColorOp{Color: color.RGBA{R: 200, G: 200, B: 200, A: 255}}.Add(gtx.Ops)
	paint.PaintOp{Rect: f32.Rectangle{Max: f32.Point{X: float32(200), Y: float32(200)}}}.Add(gtx.Ops)
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
			draw.Draw(labImg, r, src, image.Point{}, draw.Src)
		}
	}

	v.labImgOp = paint.NewImageOp(labImg)

	v.gameCounter = m.Counter
}
