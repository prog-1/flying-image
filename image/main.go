package main

import (
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	screenWidth  = 640
	screenHeight = 480
	speed        = 0.4
)

type Point struct {
	x, y float64
}

type Image struct {
	image *ebiten.Image

	pos   Point
	vel   Point
	count int
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.image.Draw(screen)
}

func (img *Image) Draw(screen *ebiten.Image) {
	w, h := img.image.Size()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	op.GeoM.Rotate(float64(int(img.count)%360) * 2 * math.Pi / 360)
	op.GeoM.Translate(img.pos.x, img.pos.y)
	screen.DrawImage(img.image, op)

}

type Game struct {
	width, height int
	last          time.Time
	image         *Image
}

func randv() (float64, float64) {
	v1 := rand.Intn(2)
	v2 := rand.Intn(2)
	if int(v1)%2 == 0 && int(v2)%2 == 0 {
		v1 = -1
		v2 = -1
		return float64(v1), float64(v2)
	} else if int(v1)%2 == 0 && int(v2)%2 != 0 {
		v1 = 1
		v2 = -1
		return float64(v1), float64(v2)
	} else if int(v1)%2 != 0 && int(v2)%2 == 0 {
		v1 = -1
		v2 = 1
		return float64(v1), float64(v2)
	} else {
		v1 = 1
		v2 = 1
		return float64(v1), float64(v2)
	}
}

func newImg(width, height int, img *ebiten.Image) *Image {
	v1, v2 := randv()
	return &Image{
		image: img,
		pos:   Point{x: float64(width / 2), y: float64(height / 2)},
		vel: Point{
			x: math.Cos(math.Pi/4) * speed * v1,
			y: math.Sin(math.Pi/4) * speed * v2,
		},
	}
}

func NewGame(width, height int, img *ebiten.Image) *Game {
	return &Game{
		width:  width,
		height: height,
		last:   time.Now(),
		image:  newImg((width)/2, (height)/2, img),
	}
}

func (g *Game) Layout(outWidth, outHeight int) (w, h int) {
	return g.width, g.height
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		os.Exit(0)
	}

	t := time.Now()
	dt := float64(t.Sub(g.last).Milliseconds())
	g.last = t
	g.image.Update(dt, g.width, g.height)
	return nil
}

func (img *Image) Update(dtMs float64, fieldWidth, fieldHeight int) {
	img.count++
	img.pos.x += img.vel.x * dtMs
	img.pos.y += img.vel.y * dtMs
	switch {
	case img.pos.x+279/2 >= float64(fieldWidth):
		img.pos.x = float64(fieldWidth) - 279/2
		img.vel.x = -img.vel.x
	case img.pos.x-279/2 <= 0:
		img.pos.x = 279 / 2
		img.vel.x = -img.vel.x
	case img.pos.y+380 >= float64(fieldHeight):
		img.pos.y = float64(fieldHeight) - 380/2
		img.vel.y = -img.vel.y
	case img.pos.y-380/2 <= 0:
		img.pos.y = 380 / 2
		img.vel.y = -img.vel.y
	}
}

func main() {
	image, _, err := ebitenutil.NewImageFromFile("kotvsp.png")
	if err != nil {
		panic(err)
	}

	rand.Seed(time.Now().UnixNano())
	g := NewGame(screenWidth, screenHeight, image)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
