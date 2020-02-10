package main

import (
	"gioui.org/app"
	"gioui.org/unit"
	"github.com/icza/golab/control"
	"github.com/icza/golab/view"
)

func main() {
	engine := control.NewEngine()

	go func() {
		w := app.NewWindow(app.Title("Gopher's Labyrinth"), app.Size(unit.Px(700), unit.Px(700)))

		v := view.New(engine, w)
		v.Loop()
	}()

	app.Main()
}
