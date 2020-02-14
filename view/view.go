// Package view is the view of the game. It handles user input and presents the game to the user.
package view

import (
	"image"
	"image/color"
	"image/draw"
	"log"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/icza/golab/engine"
)

const (
	controlsHeight = 50
	viewWidth      = 700
	viewHeight     = 700
	// WindowWidth is the suggested window width
	WindowWidth = viewWidth
	// WindowHeight is the suggested window height
	WindowHeight = controlsHeight + viewHeight
)

func init() {
	gofont.Register()
}

// imageOp wraps a paint.ImageOp and the source image.
type imageOp struct {
	paint.ImageOp
	src image.Image
}

func newImageOp(src image.Image) imageOp {
	return imageOp{
		ImageOp: paint.NewImageOp(src),
		src:     src,
	}
}

// View displays the game window and handles user input.
type View struct {
	engine *engine.Engine
	w      *app.Window
	th     *material.Theme
	gtx    *layout.Context

	// New Game button model
	newGameBtn *widget.Button

	// Difficulty options
	diffOpt *options
	// Lab size options
	labSizeOpt *options
	// Speed options
	speedOpt *options

	// "static" imageOps
	imgOpGophers  []imageOp
	imgOpDead     imageOp
	imgOpBulldogs []imageOp
	imgOpMarker   imageOp
	imgOpExit     imageOp
	imgOpWon      imageOp

	// gameCounter for the cached data
	gameCounter int
	// cached ImageOp of the whole labyrinth (only the blocks)
	labImgOp imageOp

	// Tells what offset was last applied to draw the lab view.
	// Used when calculating click position in the lab.
	labViewOffset f32.Point
	// labViewClip tells what clip rectangle was applied to draw the lab view.
	// Used to tell if a click is accepted in the lab.
	labViewClip f32.Rectangle
}

// New returns a new View.
func New(eng *engine.Engine, w *app.Window) *View {
	v := &View{
		engine:      eng,
		w:           w,
		th:          material.NewTheme(),
		gtx:         layout.NewContext((w.Queue())),
		newGameBtn:  new(widget.Button),
		imgOpDead:   newImageOp(imgDead),
		imgOpMarker: newImageOp(imgMarker),
		imgOpExit:   newImageOp(imgExit),
		imgOpWon:    newImageOp(imgWon),
	}

	for _, img := range imgGophers {
		v.imgOpGophers = append(v.imgOpGophers, newImageOp(img))
	}
	for _, img := range imgBulldogs {
		v.imgOpBulldogs = append(v.imgOpBulldogs, newImageOp(img))
	}

	v.diffOpt = newOptions(v, "Difficulty", engine.Difficulties, engine.DifficultyDefaultIdx)
	v.labSizeOpt = newOptions(v, "Lab size", engine.LabSizes, engine.LabSizeDefaultIdx)
	v.speedOpt = newOptions(v, "Speed", engine.Speeds, engine.SpeedDefaultIdx)

	return v
}

// Loop starts handing user input and frame redraws.
// This function returns only if the user closes the app.
func (v *View) Loop() {
	for e := range v.w.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			v.drawFrame(e)
		case pointer.Event:
			// TODO maybe send click on Release?
			if e.Type == pointer.Press {
				pos := e.Position.Sub(v.labViewOffset)
				// apply clip
				r := v.labViewClip
				if pos.X >= r.Min.X && pos.X < r.Max.X &&
					pos.Y >= r.Min.Y && pos.Y < r.Max.Y {
					// TODO if e.Source == pointer.Touch, set left button?
					v.engine.SendClick(engine.Click{
						X:     int(pos.X),
						Y:     int(pos.Y),
						Left:  e.Buttons&pointer.ButtonLeft != 0,
						Right: e.Buttons&pointer.ButtonRight != 0,
					})
				}
			}
		case key.Event:
			sendKey := func(dir engine.Dir) {
				v.engine.SendKey(engine.Key{DirKeys: map[engine.Dir]bool{dir: true}})
			}
			switch e.Name {
			case key.NameLeftArrow:
				sendKey(engine.DirLeft)
			case key.NameRightArrow:
				sendKey(engine.DirRight)
			case key.NameUpArrow:
				sendKey(engine.DirUp)
			case key.NameDownArrow:
				sendKey(engine.DirDown)
			}
		case system.DestroyEvent:
			log.Println("Goodbye!")
		}
	}
}

// drawFrame draws a frame of the window.
func (v *View) drawFrame(e system.FrameEvent) {
	gtx := v.gtx

	gtx.Reset(e.Config, e.Size)

	// Handle button clicks
	for v.newGameBtn.Clicked(v.gtx) {
		v.engine.NewGame(engine.GameConfig{
			Difficulty: v.diffOpt.selected().(*engine.Difficulty),
			LabSize:    v.labSizeOpt.selected().(*engine.LabSize),
			Speed:      v.speedOpt.selected().(*engine.Speed),
		})
	}
	v.diffOpt.handleInput()
	v.labSizeOpt.handleInput()
	v.speedOpt.handleInput()

	v.drawControls()
	v.drawLab()

	e.Frame(gtx.Ops)
}

// drawControls draws the control and setup widgets.
func (v *View) drawControls() {
	th, gtx := v.th, v.gtx

	layout.N.Layout(gtx, func() {
		layout.UniformInset(unit.Px(5)).Layout(gtx, func() {
			layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func() {
					layout.Inset{Left: unit.Px(10), Right: unit.Px(10)}.Layout(gtx, func() {
						b := th.Button("New Game")
						b.Background = color.RGBA{R: 20, G: 130, B: 20, A: 255}
						b.Layout(gtx, v.newGameBtn)
					})
				}),
				layout.Rigid(v.diffOpt.layout),
				layout.Rigid(v.labSizeOpt.layout),
				layout.Rigid(v.speedOpt.layout),
			)
		})
	})
}

// drawLab draws the labyrinth.
func (v *View) drawLab() {
	m := v.engine.Model
	m.RLock()
	defer m.RUnlock()

	gtx := v.gtx

	// Victory image must be drawn while locking but
	// after transformations undone.
	defer func() {
		if m.Won {
			v.drawImg(v.imgOpWon,
				v.labViewOffset.X+v.labViewClip.Min.X+(v.labViewClip.Dx()-float32(v.imgOpWon.src.Bounds().Dx()))/2,
				v.labViewOffset.Y+v.labViewClip.Min.Y+(v.labViewClip.Dy()-float32(v.imgOpWon.src.Bounds().Dy()))/2,
			)
		}
	}()

	var stack op.StackOp
	stack.Push(gtx.Ops)
	defer stack.Pop()

	// Center lab view in window:
	displayWidth, displayHeight := float32(viewWidth), float32(viewHeight)
	labWidth := float32(m.Cols * engine.BlockSize)
	labHeight := float32(m.Rows * engine.BlockSize)
	if labWidth < displayWidth {
		displayWidth = labWidth
	}
	if labHeight < displayHeight {
		displayHeight = labHeight
	}

	// Calculate the visible window of the lab image.
	// Try to center Gopher in view:
	rect := f32.Rectangle{}
	rect.Max = f32.Point{X: displayWidth, Y: displayHeight}
	rect = rect.Add(f32.Point{
		X: float32(m.Gopher.Pos.X) - displayWidth/2,
		Y: float32(m.Gopher.Pos.Y) - displayHeight/2,
	})
	// But needs correction at the edges of the view (it can't be centered)
	corr := f32.Point{}
	if rect.Min.X < 0 {
		corr.X = -rect.Min.X
	}
	if rect.Min.Y < 0 {
		corr.Y = -rect.Min.Y
	}
	if rect.Max.X > labWidth {
		corr.X = labWidth - rect.Max.X
	}
	if rect.Max.Y > labHeight {
		corr.Y = labHeight - rect.Max.Y
	}
	rect = rect.Add(corr)

	v.labViewOffset = f32.Point{
		X: (float32(gtx.Constraints.Width.Max) - displayWidth) / 2,
		Y: controlsHeight,
	}.Sub(rect.Min)
	op.TransformOp{}.Offset(v.labViewOffset).Add(gtx.Ops)
	v.labViewClip = rect
	clip.Rect{Rect: v.labViewClip}.Op(gtx.Ops).Add(gtx.Ops)

	// First the blocks:
	v.ensureLabImgOp()
	v.drawImg(v.labImgOp, 0, 0)

	// Now objects in the lab:
	// TODO do not draw images outside of the view

	// Draw target position markers:
	mbounds := imgMarker.Bounds()
	tp := m.Gopher.TargetPos
	v.drawImg(v.imgOpMarker, float32(tp.X-mbounds.Dx()/2), float32(tp.Y-mbounds.Dy()/2))
	for _, tp := range m.TargetPoss {
		v.drawImg(v.imgOpMarker, float32(tp.X-mbounds.Dx()/2), float32(tp.Y-mbounds.Dy()/2))
	}
	// Gopher:
	if m.Dead {
		v.drawObj(v.imgOpDead, m.Gopher)
	} else {
		v.drawObj(v.imgOpGophers[m.Gopher.Dir], m.Gopher)
	}
	// Bulldogs:
	for _, bd := range m.Bulldogs {
		v.drawObj(v.imgOpBulldogs[bd.Dir], bd)
	}
}

// drawObj draws the given image of the given moving obj.
func (v *View) drawObj(iop imageOp, obj *engine.MovingObj) {
	v.drawImg(iop, float32(obj.Pos.X-engine.BlockSize/2), float32(obj.Pos.Y-engine.BlockSize/2))
}

// drawImg draws the given image to the given position.
func (v *View) drawImg(iop imageOp, x, y float32) {
	var stack op.StackOp
	stack.Push(v.gtx.Ops)

	op.TransformOp{}.Offset(f32.Point{X: x, Y: y}).Add(v.gtx.Ops)

	iop.Add(v.gtx.Ops)
	imgBounds := iop.src.Bounds()
	paint.PaintOp{Rect: f32.Rectangle{
		Max: f32.Point{X: float32(imgBounds.Max.X), Y: float32(imgBounds.Max.Y)},
	}}.Add(v.gtx.Ops)

	stack.Pop()
}

// ensureLabImgOp ensures labImgOp is up-to-date
func (v *View) ensureLabImgOp() {
	m := v.engine.Model
	if v.gameCounter == m.Counter {
		// We have the lab image of the current game
		return
	}

	labImg := image.NewRGBA(image.Rect(0, 0, m.Cols*engine.BlockSize, m.Rows*engine.BlockSize))
	var r image.Rectangle
	for row := range m.Lab {
		r.Min.Y = row * engine.BlockSize
		r.Max.Y = r.Min.Y + engine.BlockSize
		for col, block := range m.Lab[row] {
			r.Min.X = col * engine.BlockSize
			r.Max.X = r.Min.X + engine.BlockSize
			src := imgBlocks[block]
			draw.Draw(labImg, r, src, image.Point{}, draw.Over)
		}
	}

	// Exit sign:
	r.Min = m.ExitPos
	r.Min = r.Min.Add(image.Point{-engine.BlockSize / 2, -engine.BlockSize / 2})
	r.Max = r.Min.Add(image.Point{engine.BlockSize, engine.BlockSize})
	draw.Draw(labImg, r, imgExit, image.Point{}, draw.Over)

	v.labImgOp = newImageOp(labImg)

	v.gameCounter = m.Counter
}
