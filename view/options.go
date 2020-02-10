package view

import (
	"fmt"
	"image/color"
	"reflect"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
)

// options groups functionality to present a widget that can be used to select one option from many.
//
// This implementation uses a Button which loops through the possible values on clicks.
type options struct {
	v      *View
	title  string
	values interface{}
	idx    int // selected index

	btn *widget.Button
}

// newOptions creates a new options.
//
// values must be a slice of possible values.
func newOptions(v *View, title string, values interface{}, defaultIdx int) *options {
	return &options{
		v:      v,
		title:  title,
		values: values,
		idx:    defaultIdx,
		btn:    new(widget.Button),
	}
}

// handleInput handles user inputs that may change the selected option.
func (o *options) handleInput() {
	for o.btn.Clicked(o.v.gtx) {
		o.idx = (o.idx + 1) % reflect.ValueOf(o.values).Len()
	}
}

// layout lays out the UI widget
func (o *options) layout() {
	layout.Inset{Left: unit.Px(10), Right: unit.Px(10)}.Layout(o.v.gtx, func() {
		b := o.v.th.Button(
			fmt.Sprintf("%s: %s", o.title, reflect.ValueOf(o.values).Index(o.idx).Interface()),
		)
		b.Background = color.RGBA{R: 100, G: 100, A: 255}
		b.Layout(o.v.gtx, o.btn)
	})
}
