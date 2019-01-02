package graphics

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Key struct {
	Up, Down, Left, Right bool
}

const (
	Width, Height int  = 768, 576
	Smooth        bool = true
	Rate          int  = 60
)

var (
	title  string     = "Canvas"
	frames int        = 0
	colour color.RGBA = color.RGBA{0xe0, 0xe0, 0xe0, 0xff}
)

func mojaveWorkaround(win *pixelgl.Window) {
	pos := win.GetPos()
	win.SetPos(pixel.ZV)
	win.SetPos(pos)
}

func displayRate(win *pixelgl.Window, cfg pixelgl.WindowConfig, s <-chan time.Time) {
	frames++
	select {
	case <-s:
		win.SetTitle(fmt.Sprintf("%s - %d fps", cfg.Title, frames))
		frames = 0
	default: // pass
	}
}

func events(win *pixelgl.Window) *Key {
	key := Key{}
	if win.Pressed(pixelgl.KeyW) {
		key.Up = true
	}
	if win.Pressed(pixelgl.KeyA) {
		key.Left = true
	}
	if win.Pressed(pixelgl.KeyS) {
		key.Down = true
	}
	if win.Pressed(pixelgl.KeyD) {
		key.Right = true
	}
	return &key
}

func run(render func(dt float64, k *Key) image.Image) {
	cfg := pixelgl.WindowConfig{
		Title:  title,
		Bounds: pixel.R(0, 0, float64(Width), float64(Height)),
		// VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	mojaveWorkaround(win)

	win.SetSmooth(Smooth)

	st := time.Tick(time.Second)
	rt := time.Tick(time.Second / time.Duration(Rate))
	ct := time.Now()

	for !win.Closed() {
		select {
		case <-rt:
			dt := time.Since(ct).Seconds()
			ct = time.Now()
			ev := events(win)

			win.Clear(colour)

			canvas := pixel.PictureDataFromImage(render(dt, ev))
			sprite := pixel.NewSprite(canvas, canvas.Bounds())
			sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))

			win.Update()

			displayRate(win, cfg, st)
		default: // pass
		}
	}
}

func Setup(render func(dt float64, k *Key) image.Image) {
	pixelgl.Run(func() { run(render) })
}
