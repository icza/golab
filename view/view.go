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
	// Difficulty button model
	diffBtn *widget.Button
	// Lab size button model
	labSizeBtn *widget.Button
	// Speed button model
	speedBtn *widget.Button
	// Selected difficulty index (in ctrl.Difficulties)
	diffIdx int
	// Selected lab size index (in ctrl.LabSizes)
	labSizeIdx int
	// Selected speed index (in ctrl.Speeds)
	speedIdx int
}

// New returns a new View.
func New(engine *ctrl.Engine, w *app.Window) *View {
	return &View{
		engine:     engine,
		w:          w,
		th:         material.NewTheme(),
		gtx:        layout.NewContext((w.Queue())),
		newGameBtn: new(widget.Button),
		diffBtn:    new(widget.Button),
		labSizeBtn: new(widget.Button),
		speedBtn:   new(widget.Button),
		diffIdx:    ctrl.DifficultyDefaultIdx,
		labSizeIdx: ctrl.DefaultLabSizeIdx,
		speedIdx:   ctrl.SpeedDefaultIdx,
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
		v.diffIdx = (v.diffIdx + 1) % len(ctrl.Difficulties)
	}
	for v.labSizeBtn.Clicked(gtx) {
		v.labSizeIdx = (v.labSizeIdx + 1) % len(ctrl.LabSizes)
	}
	for v.speedBtn.Clicked(gtx) {
		v.speedIdx = (v.speedIdx + 1) % len(ctrl.Speeds)
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
						th.Button("Difficulty: "+ctrl.Difficulties[v.diffIdx].String()).Layout(gtx, v.diffBtn)
					})
				}),
				layout.Rigid(func() {
					layout.Inset{Left: unit.Px(10), Right: unit.Px(10)}.Layout(gtx, func() {
						th.Button("Lab size: "+ctrl.LabSizes[v.labSizeIdx].String()).Layout(gtx, v.labSizeBtn)
					})
				}),
				layout.Rigid(func() {
					layout.Inset{Left: unit.Px(10), Right: unit.Px(10)}.Layout(gtx, func() {
						th.Button("Speed: "+ctrl.Speeds[v.speedIdx].String()).Layout(gtx, v.speedBtn)
					})
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
