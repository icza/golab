package main

import (
	"gioui.org/app"
	"gioui.org/unit"
	"github.com/icza/golab/engine"
	"github.com/icza/golab/view"
)

func main() {
	go func() {
		w := app.NewWindow(
			app.Title("Gopher's Labyrinth"),
			app.Size(unit.Px(view.WindowWidthPx), unit.Px(view.WindowHeightPx)),
		)

		eng := engine.NewEngine(w.Invalidate)
		go eng.Loop()

		v := view.New(eng, w)
		v.Loop()
	}()

	app.Main()
}
