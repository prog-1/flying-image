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
)

type Point struct {
	x, y float64
}

type img struct {
	img *ebiten.Image
	pos Point
	vel Point
}

type Game struct {
	width, height int
	img           *img
	last          time.Time
}

func startDir() float64 {
	a := rand.Float64()
	b := rand.Intn(2)
	if b == 0 {
		return -a
	}
	return a
}
func NewImage(width, height int, image *ebiten.Image) *img {
	speed := startDir()

	return &img{
		img: image,
		pos: Point{x: float64(width), y: float64(height)},
		vel: Point{
			x: math.Cos(math.Pi/4) * speed,
			y: math.Sin(math.Pi/4) * speed,
		},
	}
}

func NewGame(width, height int, image *ebiten.Image) *Game {
	return &Game{
		width:  width,
		height: height,
		img:    NewImage(width, height, image),
		last:   time.Now(),
	}
}

func (img *img) Update(dtMs float64, fieldWidth, fieldHeight int) {
	img.pos.x += img.vel.x * dtMs
	img.pos.y += img.vel.y * dtMs
	w, h := img.img.Size()
	if img.pos.x < 0 {
		img.pos.x = 0
		img.vel.x = -img.vel.x
	}
	if img.pos.x > float64(fieldWidth-w) {
		img.pos.x = float64(fieldWidth - w)
		img.vel.x = -img.vel.x
	}
	if img.pos.y < 0 {
		img.pos.y = 0
		img.vel.y = -img.vel.y
	}
	if img.pos.y > float64(fieldHeight-h) {
		img.pos.y = float64(fieldHeight - h)
		img.vel.y = -img.vel.y
	}
}

func (img *img) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(img.pos.x, img.pos.y)
	screen.DrawImage(img.img, op)
}

func (g *Game) Update() error {
	t := time.Now()
	dt := float64(t.Sub(g.last).Milliseconds())
	g.last = t
	g.img.Update(dt, g.width, g.height)
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	g.img.Draw(screen)
}
func (g *Game) Layout(outWidth, outHeight int) (w, h int) {
	return g.width, g.height
}

func main() {
	rand.Seed(time.Now().UnixNano())
	image, _, err := ebitenutil.NewImageFromFile("sobakentr.png")
	if err != nil || image == nil {
		log.Fatal(err)
	}
	g := NewGame(screenWidth, screenHeight, image)
	ebiten.SetWindowTitle("sabaka")
	ebiten.SetWindowSize(screenWidth, screenHeight)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
