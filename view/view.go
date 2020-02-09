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
	"github.com/icza/golab/engine"
)

func init() {
	gofont.Register()
}

// View displays the game window and handles user input.
type View struct {
	w   *app.Window
	th  *material.Theme
	gtx *layout.Context

	// New Game button model
	newGameBtn *widget.Button
	// Difficulty button model
	diffBtn *widget.Button
	// Lab size button model
	labSizeBtn *widget.Button
	// Selected difficulty index (in engine.Difficulties)
	diffIdx int
	// Selected lab size index (in engine.LabSizes)
	labSizeIdx int
}

// New returns a new View.
func New(w *app.Window) *View {
	return &View{
		w:          w,
		th:         material.NewTheme(),
		gtx:        layout.NewContext((w.Queue())),
		newGameBtn: new(widget.Button),
		diffBtn:    new(widget.Button),
		labSizeBtn: new(widget.Button),
		diffIdx:    engine.DifficultyDefaultIdx,
		labSizeIdx: engine.DefaultLabSizeIdx,
	}
}

// Loop starts handing user input. This function only returns if the user closes the app.
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

	th := v.th
	gtx := v.gtx

	gtx.Reset(e.Config, e.Size)

	for v.diffBtn.Clicked(gtx) {
		v.diffIdx = (v.diffIdx + 1) % len(engine.Difficulties)
	}
	for v.labSizeBtn.Clicked(gtx) {
		v.labSizeIdx = (v.labSizeIdx + 1) % len(engine.LabSizes)
	}

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
						th.Button("Difficulty: "+engine.Difficulties[v.diffIdx].String()).Layout(gtx, v.diffBtn)
					})
				}),
				layout.Rigid(func() {
					layout.Inset{Left: unit.Px(10), Right: unit.Px(10)}.Layout(gtx, func() {
						th.Button("Lab size: "+engine.LabSizes[v.labSizeIdx].String()).Layout(gtx, v.labSizeBtn)
					})
				}),
			)
		})
	})

	e.Frame(gtx.Ops)
}
