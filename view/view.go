package view

import (
	"log"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
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

	for e := range w.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			log.Println("frame")
			gtx.Reset(e.Config, e.Size)

			_ = th

			e.Frame(gtx.Ops)
		case system.DestroyEvent:
			log.Println("Goodbye!")
		}
	}
}
