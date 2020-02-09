package main

import (
	"gioui.org/app"
	"gioui.org/unit"
	"github.com/icza/golab/view"
)

func main() {
	go func() {
		w := app.NewWindow(app.Title("Gopher's Labyrinth"), app.Size(unit.Px(700), unit.Px(700)))

		v := view.New(w)
		v.Loop()
	}()

	app.Main()
}
