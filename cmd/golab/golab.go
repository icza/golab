package main

import (
	"log"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

func main() {
	go func() {
		w := app.NewWindow(app.Title("Gopher's Labyrinth"), app.Size(unit.Px(700), unit.Px(700)))
		loop(w)
	}()

	app.Main()
}

func loop(w *app.Window) {
	gofont.Register()
	th := material.NewTheme()
	gtx := layout.NewContext(w.Queue())

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
