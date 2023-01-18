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
	screenWidth  = 1290
	screenHeight = 960
)

type Point struct {
	x, y float64
}

type Image struct {
	size           Point
	pos            Point
	vel            Point
	png            image.Image
	r              *ebiten.Image
	rotationi      int
	trackrotationi []int
	track          []Point
}

type Game struct {
	width, height int
	image         *Image
	last          time.Time
}

func (b *Image) Update(dtMs float64, fieldWidth, fieldHeight int) {
	b.pos.x += b.vel.x * dtMs
	b.pos.y += b.vel.y * dtMs
	switch {
	case b.pos.x >= float64(fieldWidth)-1:
		b.vel.x = -b.vel.x
	case b.pos.x <= 1:
		b.vel.x = -b.vel.x
	case b.pos.y >= float64(fieldHeight)-1:
		b.vel.y = -b.vel.y
	case b.pos.y <= 1:
		b.vel.y = -b.vel.y
	}
}

func newImage(width, height int) *Image {
	xsize, ysize := 1.5, 1.5
	xpos, ypos := width/2, height/2
	rect, png, err := ebitenutil.NewImageFromFile("Arrow.png")
	if err != nil {
		log.Fatal(err)
	}
	return &Image{Point{float64(xsize), float64(ysize)}, Point{float64(xpos), float64(ypos)}, Point{
		x: math.Cos(math.Pi/4) * 0.8,
		y: math.Sin(math.Pi/4) * 0.8,
	}, png, rect, 0, nil, nil,
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
	if len(g.image.track) < 25 {
		g.image.track = append(g.image.track, g.image.pos)
	} else {
		g.image.track = append(g.image.track[1:], g.image.pos)
	}
	if len(g.image.trackrotationi) < 25 {
		g.image.trackrotationi = append(g.image.trackrotationi, g.image.rotationi)
	} else {
		g.image.trackrotationi = append(g.image.trackrotationi[1:], g.image.rotationi)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for i, r := range g.image.track {
		w, h := g.image.r.Size()
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
		op.GeoM.Scale(1/g.image.size.x, 1/g.image.size.y)
		op.GeoM.Rotate(float64(g.image.trackrotationi[i]) * math.Pi / 360)
		op.ColorM.Scale(255, 0, 0, float64(2000)/10000)
		op.GeoM.Translate(r.x, r.y)
		screen.DrawImage(g.image.r, op)

	}
	w, h := g.image.r.Size()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	op.GeoM.Scale(1/g.image.size.x, 1/g.image.size.y)
	if g.image.vel.x > 0 && g.image.vel.y > 0 {
		g.image.rotationi = 180
	}
	if g.image.vel.x < 0 && g.image.vel.y > 0 {
		g.image.rotationi = 360
	}
	if g.image.vel.x < 0 && g.image.vel.y < 0 {
		g.image.rotationi = -180
	}
	if g.image.vel.x > 0 && g.image.vel.y < 0 {
		g.image.rotationi = 0
	}
	op.GeoM.Rotate(float64(g.image.rotationi) * math.Pi / 360)
	op.GeoM.Translate(g.image.pos.x, g.image.pos.y)
	screen.DrawImage(g.image.r, op)
}

func NewGame(width, height int) *Game {
	return &Game{
		width:  width,
		height: height,
		image:  newImage(width-1, height-1),
		last:   time.Now(),
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	g := NewGame(screenWidth, screenHeight)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
