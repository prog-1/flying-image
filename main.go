package main

import (
	"image"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 800
	screenHeight = 600
	radius       = 20
)

type Point struct {
	x, y float64
}

type myImage struct {
	i     *ebiten.Image
	img   image.Image
	pos   Point
	track []Point

	rottation    int
	vel          Point
	SizeX, SizeY float64
}

func sign() float64 {
	if rand.Intn(2) == 0 {
		return -1.0
	}
	return 1.0
}
func newImage(width, height int) *myImage {
	x0, y0 := 100, 100
	x1, y1 := 200, 200
	rect, img, err := ebitenutil.NewImageFromFile("g.png")
	if err != nil {
		log.Fatal(err)
	}
	x, y := ebiten.CursorPosition()
	return &myImage{rect, img, Point{float64(x), float64(y)}, nil, 0, Point{
		x: math.Cos(math.Pi/4) * rand.Float64() * sign(),
		y: math.Sin(math.Pi/4) * rand.Float64() * sign(),
	}, float64(x1 - x0), float64(y1 - y0)}
}

func (b *myImage) Update(dtMs float64, fieldWidth, fieldHeight int) {
	b.pos.x += b.vel.x * dtMs
	b.pos.y += b.vel.y * dtMs
	switch {
	case b.pos.x+radius >= float64(fieldWidth):
		b.pos.x = float64(fieldWidth-1) - radius
		b.vel.x = -b.vel.x
		b.vel.x *= 0.9
		b.vel.y *= 0.9
	case b.pos.x-radius <= 0:
		b.pos.x = 1 + radius
		b.vel.x = -b.vel.x
		b.vel.x *= 0.9
		b.vel.y *= 0.9
	case b.pos.y+radius >= float64(fieldHeight):
		b.pos.y = float64(fieldHeight-1) - radius
		b.vel.y = -b.vel.y
		b.vel.x *= 0.9
		b.vel.y *= 0.9
	case b.pos.y-radius <= 0:
		b.pos.y = 1 + radius
		b.vel.y = -b.vel.y
		b.vel.x *= 0.9
		b.vel.y *= 0.9
	}
}

type Game struct {
	width, height int
	image         *myImage
	last          time.Time
}

// NewGame returns a new Game instance.
func NewGame(width, height int) *Game {
	return &Game{
		width:  width,
		height: height,
		image:  newImage(width/2, height/2),
		last:   time.Now(),
	}
}

func (g *Game) Layout(outWidth, outHeight int) (w, h int) {
	return g.width, g.height
}

func (g *Game) Update() error {
	t := time.Now()
	dt := float64(t.Sub(g.last).Milliseconds())
	g.last = t
	g.image.Update(dt, g.width, g.height)
	if len(g.image.track) < 50 {
		g.image.track = append(g.image.track, g.image.pos)
	} else {
		g.image.track = append(g.image.track[1:], g.image.pos)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// tmp := g.image.img
	// trnsp := 0xff

	for i, r := range g.image.track {

		if i%2 == 0 {
			// mask := image.NewUniform(color.Alpha{128})
			// draw.DrawMask()
			w, h := g.image.i.Size()
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
			op.GeoM.Scale(50/g.image.SizeX, 50/g.image.SizeY)
			op.GeoM.Rotate(float64(int(g.image.rottation)%360) * 2 * math.Pi / 360)
			// op.GeoM.
			op.GeoM.Translate(r.x, r.y)
			screen.DrawImage(g.image.i, op)
		}
	}
	w, h := g.image.i.Size()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	op.GeoM.Scale(50/g.image.SizeX, 50/g.image.SizeY)
	op.GeoM.Rotate(float64(int(g.image.rottation)%360) * 2 * math.Pi / 360)
	op.GeoM.Translate(g.image.pos.x, g.image.pos.y)
	screen.DrawImage(g.image.i, op)

}

func main() {
	rand.Seed(time.Now().UnixNano())
	g := NewGame(screenWidth, screenHeight)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
