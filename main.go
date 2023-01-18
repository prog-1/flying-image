package main

import (
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 640
	screenHeight = 480

	// Ball radius.
	radius = 1
	// Ball default speed in px/ms.
	//speed = 0.4
)

// Point is a struct for representing 2D vectors.
type Point struct {
	x, y float64
}

type Hippo struct {
	pos Point
	// Ball speed in px/ms.
	track []Point
	vel   Point
	color color.RGBA
}

var (
	img *ebiten.Image
)

// NewBall initializes and returns a new Ball instance.
func NewHippo(x, y int) *Hippo {
	return &Hippo{
		pos: Point{x: float64(x), y: float64(y)},
		vel: Point{
			x: math.Cos(math.Pi/4) * rand.Float64(),
			y: math.Sin(math.Pi/4) * rand.Float64(),
		},
		color: color.RGBA{
			R: uint8(139),
			G: uint8(69),
			B: uint8(19),
			A: 255,
		},
	}

}

// Update changes a ball state.
//
// dtMs defines a time interval in microseconds between now and a previous time
// when Update was called.
func (b *Hippo) Update(dtMs float64, fieldWidth, fieldHeight int) {
	b.pos.x += b.vel.x * dtMs
	b.pos.y += b.vel.y * dtMs
	switch {
	case b.pos.x+radius >= float64(fieldWidth):
		b.pos.x = float64(fieldWidth-1) - radius
		b.vel.x = -b.vel.x
	case b.pos.x-radius <= 0:
		b.pos.x = 1 + radius
		b.vel.x = -b.vel.x
	case b.pos.y+radius >= float64(fieldHeight):
		b.pos.y = float64(fieldHeight-1) - radius
		b.vel.y = -b.vel.y
	case b.pos.y-radius <= 0:
		b.pos.y = 1 + radius
		b.vel.y = -b.vel.y
	}
}

// Draw renders a ball on a screen.
func (b *Hippo) Draw(screen *ebiten.Image) {

	tr := b.color
	for i := len(b.track) - 1; i >= 0; i-- {
		ebitenutil.DrawCircle(screen, b.track[i].x, b.track[i].y, 5, tr)
		tr.A -= 2
	}
	//imgCX := float64(b.img.Bounds().Dx()) / 2
	//imgCY := float64(b.img.Bounds().Dy()) / 2
	var op ebiten.DrawImageOptions
	op.GeoM.Translate(b.pos.x, b.pos.y)
	screen.DrawImage(img, &op)
}

// Game is a game instance.
type Game struct {
	width, height int
	hip           []*Hippo
	// last is a timestamp when Update was called last time.
	last time.Time
}

// NewGame returns a new Game instance.
func NewGame(width, height int) *Game {
	return &Game{
		width:  width,
		height: height,
		hip:    []*Hippo{},
		// A new ball is created at the center of the screen.
		// ball: NewBall(width/2, height/2),
		last: time.Now(),
	}
}

func (g *Game) Layout(outWidth, outHeight int) (w, h int) {
	return g.width, g.height
}

// Update updates a game state.
func (g *Game) Update() error {
	t := time.Now()
	dt := float64(t.Sub(g.last).Milliseconds())
	g.last = t
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.hip = append(g.hip, NewHippo(ebiten.CursorPosition()))
	}
	for i := range g.hip {
		g.hip[i].Update(dt, g.width, g.height)
		if len(g.hip[i].track) < 100 {
			g.hip[i].track = append(g.hip[i].track, g.hip[i].pos)
		} else {
			g.hip[i].track = append(g.hip[i].track[1:], g.hip[i].pos)
		}

	}
	return nil
}

// Draw renders a game screen.
func (g *Game) Draw(screen *ebiten.Image) {
	// w1, h1 := img.Size()
	// var op ebiten.DrawImageOptions
	// op.GeoM.Translate(float64(g.width-w1/2), float64(g.height-h1/2))

	screen.Fill(color.White)
	for i := range g.hip {
		g.hip[i].Draw(screen)
	}
}

func main() {
	var err error
	img, _, err = ebitenutil.NewImageFromFile("hipp.png")
	if err != nil {
		log.Fatal(err)
	}
	rand.Seed(time.Now().UnixNano())
	g := NewGame(screenWidth, screenHeight)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
