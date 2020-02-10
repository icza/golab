// Package view is the view of the game. It handles user input and presents the game to the user.
package view

import (
	"log"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
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
	v.labSizeOpt = newOptions(v, "Lab size", ctrl.LabSizes, ctrl.DefaultLabSizeIdx)
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
	log.Println("frame")

	gtx := v.gtx

	gtx.Reset(e.Config, e.Size)

	// Handle button clicks
	v.diffOpt.handleInput()
	v.labSizeOpt.handleInput()
	v.speedOpt.handleInput()

	v.drawControls()
	v.drawLab()

	e.Frame(gtx.Ops)
}

// drawControls draws the control and setup widgets.
func (v *View) drawControls() {
	th := v.th
	gtx := v.gtx

	layout.N.Layout(gtx, func() {
		layout.UniformInset(unit.Px(5)).Layout(gtx, func() {
			layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func() {
					layout.Inset{Left: unit.Px(10), Right: unit.Px(10)}.Layout(gtx, func() {
						th.Button("New Game").Layout(gtx, v.newGameBtn)
					})
				}),
				layout.Rigid(func() {
					v.diffOpt.layout()
				}),
				layout.Rigid(func() {
					v.labSizeOpt.layout()
				}),
				layout.Rigid(func() {
					v.speedOpt.layout()
				}),
			)
		})
	})
}

// drawLab draws the labyrinth.
func (v *View) drawLab() {
	m := v.engine.Model
	m.Mu.RLock()
	defer m.Mu.RUnlock()

	// TODO
}
