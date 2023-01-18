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

	// Ball default speed in px/ms.
	speed = 0.4
)

// Point is a struct for representing 2D vectors.
type Point struct {
	x, y float64
}

type trace struct {
	pos  Point
	time time.Time
}

type Image struct {
	image  *ebiten.Image
	traces []trace
	pos    Point
	vel    Point
}

func (img *Image) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(img.pos.x, img.pos.y)
	for i := range img.traces {
		t := time.Now()
		a := 10000 - t.Sub(img.traces[i].time).Milliseconds()
		if a > 0 {
			op2 := &ebiten.DrawImageOptions{}
			op2.GeoM.Translate(img.traces[i].pos.x, img.traces[i].pos.y)
			op2.ColorM.Scale(1, 1, 1, float64(a)/10000)
			screen.DrawImage(img.image, op2)
		}
	}
	screen.DrawImage(img.image, op)
}

func (img *Image) Update(dtMs float64, fieldWidth, fieldHeight int) {
	img.pos.x += img.vel.x * dtMs * speed
	img.pos.y += img.vel.y * dtMs * speed
	x, y := img.image.Size()
	switch {
	case x+int(img.pos.x) >= fieldWidth && img.vel.x > 0:
		img.vel.x *= -1
	case y+int(img.pos.y) >= fieldHeight && img.vel.y > 0:
		img.vel.y *= -1
	case img.pos.x <= 0 && img.vel.x < 0:
		img.vel.x *= -1
	case img.pos.y <= 0 && img.vel.y < 0:
		img.vel.y *= -1
	}
	img.traces = append(img.traces, trace{pos: Point{img.pos.x, img.pos.y}, time: time.Now()})
}

/*
func (b *Ball) Update(dtMs float64, fieldWidth, fieldHeight int) {
	// if b.
	if b.curentSpeed > 0 {
		b.pos.x += b.vel.x * dtMs * b.curentSpeed
		b.pos.y += b.vel.y * dtMs * b.curentSpeed
		// fmt.Println(b.curentSpeed)
		b.curentSpeed -= 0.01 * dtMs / 1000
	}
	switch {
	case b.pos.x+radius >= float64(fieldWidth) && b.vel.x > 0:
		b.vel.x *= -1
		//b.vel.y = math.Sin(-math.Pi/4) * speed
	case b.pos.x-radius < 0 && b.vel.x < 0:
		b.vel.x *= -1

	case b.pos.y+radius >= float64(fieldHeight) && b.vel.y > 0:
		b.vel.y *= -1

	case b.pos.y-radius < 0 && b.vel.y < 0:
		b.vel.y *= -1

	}
}
*/
// Draw renders a ball on a screen.

// Game is a game instance.
type Game struct {
	width, height int
	images        []*Image
	// last is a timestamp when Update was called last time.
	last  time.Time
	image *ebiten.Image
}

// NewGame returns a new Game instance.
func NewGame(width, height int, img *ebiten.Image) *Game {
	return &Game{
		width:  width,
		height: height,
		// A new ball is created at the center of the screen.
		images: []*Image{&Image{image: img,
			pos: Point{x: float64(width / 2), y: float64(height / 2)},
			vel: Point{
				x: math.Cos(math.Pi/4) * speed,
				y: math.Sin(math.Pi/4) * speed,
			},
		}},
		last:  time.Now(),
		image: img,
	}
}

func (g *Game) Layout(outWidth, outHeight int) (w, h int) {
	return g.width, g.height
}

/*
// Update updates a game state.
func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		g.ball = append(g.ball, NewBall(x, y))
	}
	for i := range g.ball {
		for j := range g.ball[i+1:] {
			j += i + 1
			if math.Sqrt(math.Pow(g.ball[i].pos.x-g.ball[j].pos.x, 2)+math.Pow(g.ball[i].pos.y-g.ball[j].pos.y, 2)) <= 2*radius {
				g.ball[i].vel, g.ball[j].vel = g.ball[j].vel, g.ball[i].vel
				g.ball[i].curentSpeed, g.ball[j].curentSpeed = g.ball[j].curentSpeed, g.ball[i].curentSpeed
			}
		}
	}
	t := time.Now()
	dt := float64(t.Sub(g.last).Milliseconds())
	g.last = t
	for i := range g.ball {
		fmt.Println(g.ball[i].curentSpeed)
		g.ball[i].Update(dt, g.width, g.height)
	}
	return nil
}
*/
// Draw renders a game screen.
func (g *Game) Update() error {
	t := time.Now()
	dt := float64(t.Sub(g.last).Milliseconds())
	g.last = t
	for i := range g.images {
		g.images[i].Update(dt, g.width, g.height)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	for i := range g.images {
		g.images[i].Draw(screen)
	}
}

func main() {
	img, _, err := ebitenutil.NewImageFromFile("a.png")
	if err != nil {
		panic(err)
	}
	rand.Seed(time.Now().UnixNano())
	g := NewGame(screenWidth, screenHeight, img)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
