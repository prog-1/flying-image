package main

import (
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 640
	screenHeight = 480
	speed        = 0.15
	radius       = 20
)

type Point struct {
	x, y float64
}

type Image struct {
	image *ebiten.Image
	pos   Point
	vel   Point
	size  Point
}

type Game struct {
	img  Image
	last time.Time
}

func (g *Game) Update() error {
	t := time.Now()
	dt := float64(t.Sub(g.last).Milliseconds())
	g.last = t
	g.img.Update(dt, screenWidth-220, screenHeight-220)
	return nil

}

func (img *Image) Update(dtMs float64, fieldWidth, fieldHeight int) {
	img.pos.x += img.vel.x * dtMs
	img.pos.y += img.vel.y * dtMs
	img.pos.x += img.vel.x * dtMs
	img.pos.y += img.vel.y * dtMs
	switch {
	case img.pos.x+radius >= float64(fieldWidth):
		img.pos.x = float64(fieldWidth) - radius
		img.vel.x = -img.vel.x
	case img.pos.x-radius <= 0:
		img.pos.x = radius
		img.vel.x = -img.vel.x
	case img.pos.y+radius >= float64(fieldHeight):
		img.pos.y = float64(fieldHeight) - radius
		img.vel.y = -img.vel.y
	case img.pos.y-radius <= 0:
		img.pos.y = radius
		img.vel.y = -img.vel.y
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.img.Draw(screen)
}

func (img *Image) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(img.size.x/100, img.size.y/100)
	op.GeoM.Translate(img.pos.x, img.pos.y)
	screen.DrawImage(img.image, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (w, h int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Bubble sort is our everything!")

	img := Image{
		pos:  Point{x: float64(screenWidth / 2), y: float64(screenHeight / 2)},
		vel:  Point{x: math.Cos(math.Pi/4) * speed, y: math.Sin(math.Pi/4) * speed},
		size: Point{x: 50, y: 50},
	}
	var err error
	img.image, _, err = ebitenutil.NewImageFromFile("Bubble.png")
	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().UnixNano())
	g := NewGame(screenWidth, screenHeight, img)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}

}

func NewGame(width, height int, img Image) *Game {
	return &Game{
		img:  img,
		last: time.Now(),
	}
}
