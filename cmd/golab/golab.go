package main

import (
	"gioui.org/app"
	"gioui.org/unit"
	"github.com/icza/golab/ctrl"
	"github.com/icza/golab/view"
)

func main() {
	go func() {
		w := app.NewWindow(
			app.Title("Gopher's Labyrinth"),
			app.Size(unit.Px(view.WindowWidth), unit.Px(view.WindowHeight)),
		)

		engine := ctrl.NewEngine(w.Invalidate)
		go engine.Loop()

		v := view.New(engine, w)
		v.Loop()
	}()

	app.Main()
}
