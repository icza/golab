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
		newGameBtn    = new(widget.Button)
		labSizeBtn    = new(widget.Button)
		difficultyBtn = new(widget.Button)
	)

	for e := range w.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			log.Println("frame")
			gtx.Reset(e.Config, e.Size)

			layout.NW.Layout(gtx, func() {
				layout.UniformInset(unit.Px(5)).Layout(gtx, func() {
					layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func() {
							th.Button("New game").Layout(gtx, newGameBtn)
						}),
						layout.Rigid(func() {
							layout.Inset{Left: unit.Px(20)}.Layout(gtx, func() {
								th.Body1("Difficulty:").Layout(gtx)
							})
						}),
						layout.Rigid(func() {
							th.Button("Medium").Layout(gtx, labSizeBtn)
						}),
						layout.Rigid(func() {
							layout.Inset{Left: unit.Px(20)}.Layout(gtx, func() {
								th.Body1("Labyrinth size:").Layout(gtx)
							})
						}),
						layout.Rigid(func() {
							th.Button("Normal").Layout(gtx, difficultyBtn)
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
