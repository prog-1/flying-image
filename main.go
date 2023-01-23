package main

import (
	_ "image/jpeg"
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
	trackLength  = 30
)

type Image struct {
	image *ebiten.Image
	pos   Point
	vel   Point
	track []*Track
}

type Point struct {
	x, y float64
}

type Track struct {
	pos        Point
	alphaScale float64
}

type Game struct {
	width, height int
	image         *Image
	last          time.Time
}

func (img *Image) Update(dtMs float64, fieldWidth, fieldHeight int) {
	img.pos.x += img.vel.x * dtMs
	img.pos.y += img.vel.y * dtMs
	w, h := img.image.Size()
	switch {
	case img.pos.x+float64(w) >= float64(fieldWidth):
		img.pos.x = float64(fieldWidth) - float64(w)
		img.vel.x = -img.vel.x
	case img.pos.x <= 0:
		img.pos.x = 0
		img.vel.x = -img.vel.x
	case img.pos.y+float64(h) >= float64(fieldHeight):
		img.pos.y = float64(fieldHeight) - float64(h)
		img.vel.y = -img.vel.y
	case img.pos.y <= 0:
		img.pos.y = 0
		img.vel.y = -img.vel.y
	}
	img.track = append(img.track, &Track{
		pos:        Point{x: img.pos.x, y: img.pos.y},
		alphaScale: 0.6,
	})
	for _, t := range img.track {
		t.alphaScale -= 0.6 / trackLength
	}
	if len(img.track) > trackLength {
		img.track = img.track[len(img.track)-trackLength:]
	}
}

func (img *Image) Draw(screen *ebiten.Image) {
	for _, t := range img.track {
		optionsTrack := &ebiten.DrawImageOptions{}
		optionsTrack.GeoM.Translate(t.pos.x, t.pos.y)
		optionsTrack.ColorM.Scale(1, 1, 1, t.alphaScale)
		screen.DrawImage(img.image, optionsTrack)
	}
	options := &ebiten.DrawImageOptions{}
	options.GeoM.Translate(img.pos.x, img.pos.y)
	screen.DrawImage(img.image, options)
}

func NewGame(width, height int, image *ebiten.Image) *Game {
	return &Game{
		width:  width,
		height: height,
		image: &Image{
			image: image,
			vel: Point{
				x: math.Cos(math.Pi/4) * rand.Float64(),
				y: math.Sin(math.Pi/4) * rand.Float64(),
			},
		},
		last: time.Now(),
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
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.image.Draw(screen)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	image, _, err := ebitenutil.NewImageFromFile("golang.jpg")
	if err != nil {
		log.Fatal(err)
	}
	g := NewGame(screenWidth, screenHeight, image)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
