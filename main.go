package main

import (
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
)

type Point struct {
	x, y float64
}

type game struct {
	last     time.Time
	image    *Image
	trailLen uint
}

type Image struct {
	image *ebiten.Image
	pos   Point
	vel   Point
	speed float64
	alfa  float64
}

func (g *game) Layout(outWidth, outHeight int) (w, h int) { return screenWidth, screenHeight }
func (g *game) Update() error {
	t := time.Now()
	dtMs := float64(t.Sub(g.last).Milliseconds())
	g.last = t
	g.image.Update(dtMs)
	return nil
}
func (g *game) Draw(screen *ebiten.Image) {
	g.image.Draw(screen)
}

func (i *Image) Update(dtMs float64) {
	i.pos.x += i.vel.x * dtMs
	i.pos.y += i.vel.y * dtMs

	halfWidth, halfHeight := float64(i.image.Bounds().Dx()/2), float64(i.image.Bounds().Dy()/2)
	switch {
	case i.pos.x+halfWidth >= float64(screenWidth) || i.pos.x-halfWidth <= 0:
		i.vel.x = -i.vel.x
	case i.pos.y+halfHeight >= float64(screenHeight) || i.pos.y-halfHeight <= 0:
		i.vel.y = -i.vel.y
	default:
		return
	}
	i.alfa -= 0.1
}

func (i *Image) Draw(screen *ebiten.Image) {
	halfWidth, halfHeight := float64(i.image.Bounds().Dx()/2), float64(i.image.Bounds().Dy()/2)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(i.pos.x-halfWidth, i.pos.y-halfHeight)
	op.ColorM.Scale(1, 1, 1, i.alfa)
	screen.DrawImage(i.image, op)
}

func NewImage(path string) (*Image, error) {
	var img *ebiten.Image
	var err error
	if img, _, err = ebitenutil.NewImageFromFile(path); err != nil {
		return nil, err
	}

	rad := rand.Float64() * 2 * math.Pi
	speed := 0.1

	return &Image{
		image: img,
		pos:   Point{screenWidth / 2, screenHeight / 2},
		vel:   Point{math.Cos(rad) * speed, math.Sin(rad) * speed},
		speed: speed,
		alfa:  1,
	}, nil
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Peace")

	var img *Image
	var err error
	if img, err = NewImage("Peace.png"); err != nil {
		log.Fatal(err)
	}
	g := game{time.Now(), img, 10}

	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}
}

// How to implement trail?
// Store last n positions of the image in the n-sized queue
// At each of these positions draw the same image, but with alfa = 255/n*i, where i is the index of the element
// At each update add new position to the front of the queue
