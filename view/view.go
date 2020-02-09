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
}

// New returns a new View.
func New(w *app.Window) *View {
	return &View{
		w:   w,
		th:  material.NewTheme(),
		gtx: layout.NewContext((w.Queue())),
	}
}

// Loop starts handing user input. This function only returns if the user closes the app.
func (v *View) Loop() {
	w := v.w
	th := v.th
	gtx := v.gtx

	var (
		newGameBtn = new(widget.Button)
		diffBtn    = new(widget.Button)
		labSizeBtn = new(widget.Button)
	)

	var (
		diffIdx    = engine.DifficultyDefaultIdx
		labSizeIdx = engine.DefaultLabSizeIdx
	)

	for e := range w.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			log.Println("frame")
			gtx.Reset(e.Config, e.Size)

			for diffBtn.Clicked(gtx) {
				diffIdx = (diffIdx + 1) % len(engine.Difficulties)
			}
			for labSizeBtn.Clicked(gtx) {
				labSizeIdx = (labSizeIdx + 1) % len(engine.LabSizes)
			}

			layout.N.Layout(gtx, func() {
				layout.UniformInset(unit.Px(5)).Layout(gtx, func() {
					layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func() {
							layout.Inset{Left: unit.Px(10), Right: unit.Px(10)}.Layout(gtx, func() {
								th.Button("New Game").Layout(gtx, newGameBtn)
							})
						}),
						layout.Rigid(func() {
							layout.Inset{Left: unit.Px(10), Right: unit.Px(10)}.Layout(gtx, func() {
								th.Button("Difficulty: "+engine.Difficulties[diffIdx].String()).Layout(gtx, diffBtn)
							})
						}),
						layout.Rigid(func() {
							layout.Inset{Left: unit.Px(10), Right: unit.Px(10)}.Layout(gtx, func() {
								th.Button("Lab size: "+engine.LabSizes[labSizeIdx].String()).Layout(gtx, labSizeBtn)
							})
						}),
					)
				})
			})

			e.Frame(gtx.Ops)
		case system.DestroyEvent:
			log.Println("Goodbye!")
		}
	}
}
