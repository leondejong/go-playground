package game

import (
	"image"
	"image/color"
	"math"

	"github.com/fogleman/gg"
	"github.com/leondejong/go-playground/graphics"
)

type rectangle struct {
	x, y, w, h float64
}

const (
	w, h, s, g, i float64 = 16, 16, 320, 16, 8
)

var (
	x, y, vx, vy float64     = 32, 32, 0, 0
	ctx          *gg.Context = gg.NewContext(graphics.Width, graphics.Height)
	message      string      = ""
)

var level = []*rectangle{
	// Border
	{0, 0, 768, 16},
	{0, 560, 768, 16},
	{0, 0, 16, 576},
	{752, 0, 16, 576},
	// floors
	{336, 144, 16, 288},
	{352, 144, 336, 16},
	{418, 236, 336, 16},
	{352, 326, 336, 16},
	{464, 416, 112, 16},
	{640, 416, 112, 16},
	{576, 486, 64, 16},
	// platforms
	{80, 486, 64, 16},
	{208, 416, 64, 16},
	{80, 348, 64, 16},
	{208, 280, 64, 16},
	{80, 212, 64, 16},
	{208, 144, 64, 16},
	// stairs
	{448, 432, 16, 16},
	{432, 448, 16, 16},
	{416, 464, 16, 16},
	{400, 480, 16, 16},
	{384, 496, 16, 16},
	{368, 512, 16, 16},
	{352, 528, 16, 16},
	{336, 544, 16, 16},
	// walls
	{420, 80, 16, 64},
	{588, 80, 16, 64},
	{504, 16, 16, 64},
}

func reset() {
	ctx = gg.NewContext(graphics.Width, graphics.Height)
}

func intersect(a, b *rectangle) bool {
	n := a.y < b.y+b.h
	s := b.y < a.y+a.h
	e := b.x < a.x+a.w
	w := a.x < b.x+b.w
	return n && s && e && w
}

func collision(a *rectangle) bool {
	c := false
	for _, b := range level {
		if intersect(a, b) {
			c = true
		}
	}
	return c
}

func drawRectangle(r *rectangle, c color.RGBA) {
	ctx.DrawRectangle(r.x, r.y, r.w, r.h)
	ctx.SetColor(c)
	ctx.Fill()
}

func drawLevel() {
	for _, rect := range level {
		drawRectangle(rect, color.RGBA{0x00, 0xe0, 0x90, 0xff})
	}
}

func translate() {
	dh := rectangle{x + vx, y, w, h}
	dv := rectangle{x, y + vy, w, h}
	for _, u := range level {
		if intersect(u, &dh) {
			if vx < 0 {
				vx = u.x + u.w - x
			} else if vx > 0 {
				vx = u.x - x - w
			}
		}
		if intersect(u, &dv) {
			if vy < 0 {
				vy = u.y + u.h - y
			} else if vy > 0 {
				vy = u.y - y - h
			}
		}
	}
	x += vx
	y += vy
}

func update(dt float64, k *graphics.Key) image.Image {
	vx = 0
	if k.Left {
		vx = -s * dt
	}
	if k.Right {
		vx = s * dt
	}
	vy += g * dt
	dx := x + vx
	dy := y + vy
	translate()
	if dx != x {
		vx = 0
	}
	if dy != y {
		vy = 0
	}
	if y < dy && k.Up {
		vy -= i
	}
	return render(dt)
}

func render(dt float64) image.Image {
	reset()
	drawLevel()
	player := rectangle{math.Round(x), math.Round(y), w, h}
	colour := color.RGBA{0x00, 0x90, 0xff, 0xff}
	if collision(&player) {
		colour = color.RGBA{0xff, 0x40, 0x00, 0xff}
	}
	drawRectangle(&player, colour)
	ctx.DrawString(message, 24, 32)
	return ctx.Image()
}

func Init() {
	graphics.Setup(update)
}
