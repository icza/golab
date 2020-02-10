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
	"github.com/icza/golab/control"
)

func init() {
	gofont.Register()
}

// View displays the game window and handles user input.
type View struct {
	engine *control.Engine
	w      *app.Window
	th     *material.Theme
	gtx    *layout.Context

	// New Game button model
	newGameBtn *widget.Button
	// Difficulty button model
	diffBtn *widget.Button
	// Lab size button model
	labSizeBtn *widget.Button
	// Selected difficulty index (in control.Difficulties)
	diffIdx int
	// Selected lab size index (in control.LabSizes)
	labSizeIdx int
}

// New returns a new View.
func New(engine *control.Engine, w *app.Window) *View {
	return &View{
		engine:     engine,
		w:          w,
		th:         material.NewTheme(),
		gtx:        layout.NewContext((w.Queue())),
		newGameBtn: new(widget.Button),
		diffBtn:    new(widget.Button),
		labSizeBtn: new(widget.Button),
		diffIdx:    control.DifficultyDefaultIdx,
		labSizeIdx: control.DefaultLabSizeIdx,
	}
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
	for v.diffBtn.Clicked(gtx) {
		v.diffIdx = (v.diffIdx + 1) % len(control.Difficulties)
	}
	for v.labSizeBtn.Clicked(gtx) {
		v.labSizeIdx = (v.labSizeIdx + 1) % len(control.LabSizes)
	}

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
					layout.Inset{Left: unit.Px(10), Right: unit.Px(10)}.Layout(gtx, func() {
						th.Button("Difficulty: "+control.Difficulties[v.diffIdx].String()).Layout(gtx, v.diffBtn)
					})
				}),
				layout.Rigid(func() {
					layout.Inset{Left: unit.Px(10), Right: unit.Px(10)}.Layout(gtx, func() {
						th.Button("Lab size: "+control.LabSizes[v.labSizeIdx].String()).Layout(gtx, v.labSizeBtn)
					})
				}),
			)
		})
	})
}

// drawLab draws the labyrinth.
func (v *View) drawLab() {
	// TODO
}
