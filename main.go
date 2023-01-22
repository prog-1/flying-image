package main

import (
	"log"
	"math"
	"math/rand"
	"strconv"
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

type List struct {
	maxLen uint
	len    uint
	head   *Node
}

func (l *List) Prepend(pos Point) {
	el := &Node{pos, nil}
	el.next = l.head
	l.head = el
	l.len++
	if l.len > l.maxLen {
		cur := l.head
		for i := uint(0); i < l.maxLen; i++ {
			cur = cur.next
		}
		cur.next = nil
	}
}

type Node struct {
	pos  Point
	next *Node
}

type game struct {
	last  time.Time
	image *Image
	trail *List
}

type Image struct {
	image *ebiten.Image
	pos   Point
	vel   Point
	speed float64
}

func (g *game) Layout(outWidth, outHeight int) (w, h int) { return screenWidth, screenHeight }
func (g *game) Update() error {
	t := time.Now()
	dtMs := float64(t.Sub(g.last).Milliseconds())
	g.last = t

	g.trail.Prepend(g.image.pos)

	g.image.Update(dtMs)
	return nil
}
func (g *game) Draw(screen *ebiten.Image) {
	halfWidth, halfHeight := float64(g.image.image.Bounds().Dx()/2), float64(g.image.image.Bounds().Dy()/2)
	for i, cur := 0, g.trail.head; cur != nil; i, cur = i+1, cur.next {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(cur.pos.x-halfWidth, cur.pos.y-halfHeight)
		alfa := float64(1) / float64(i+1)
		ebitenutil.DebugPrint(screen, strconv.FormatFloat(alfa, 'f', -1, 32))
		op.ColorM.Scale(1, 1, 1, alfa)
		screen.DrawImage(g.image.image, op)
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(g.image.pos.x-halfWidth, g.image.pos.y-halfHeight)
	screen.DrawImage(g.image.image, op)
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
	g := game{time.Now(), img, &List{maxLen: 100}}

	if err := ebiten.RunGame(&g); err != nil {
		log.Fatal(err)
	}
}

// How to implement trail?
// Store last n positions of the image in the n-sized queue
// At each update add new position to the front of the queue
// At each of these positions draw the same image, but with alfa = 255/n*i, where i is the index of the element
//alfa := float64(1) / float64(g.trail.len) * float64(i+1)
// No, it doesn't work as expected..
