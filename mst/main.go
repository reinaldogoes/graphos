package main

import (
	"image"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
)

const (
	border       = 10
	screenWidth  = 320 + border*2
	screenHeight = 240 + border*2
)

var (
	tmpImg       *ebiten.Image
	img          *image.RGBA
	uTime        uint64
	updateScreen bool
)

var CGAColors = []struct {
	R byte
	G byte
	B byte
}{
	{0, 0, 0},
	{0, 0, 170},
	{0, 170, 0},
	{0, 170, 170},
	{170, 0, 0},
	{170, 0, 170},
	{170, 85, 0},
	{170, 170, 170},
	{85, 85, 85},
	{85, 85, 255},
	{85, 255, 85},
	{85, 255, 255},
	{255, 85, 85},
	{255, 85, 255},
	{255, 255, 85},
	{255, 255, 255},
}

var lastKey = struct {
	Time uint64
	Char byte
}{
	0,
	0,
}

func mergeColorCode(b, f byte) byte {
	return (f & 0xff) | (b << 4)
}

func drawPix(x, y int, color byte) {
	x += border
	y += border
	if x < border || y < border || x >= screenWidth-border || y >= screenHeight-border {
		return
	}
	pos := 4*y*screenWidth + 4*x
	img.Pix[pos] = CGAColors[color].R
	img.Pix[pos+1] = CGAColors[color].G
	img.Pix[pos+2] = CGAColors[color].B
	img.Pix[pos+3] = 0xff
	updateScreen = true
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func update(screen *ebiten.Image) error {

	uTime++

	if updateScreen {
		tmpImg.ReplacePixels(img.Pix)
		updateScreen = false
	}

	screen.DrawImage(tmpImg, nil)
	input()
	return nil
}

func main() {
	rand.Seed(time.Now().Unix())

	for i := 0; i < 10; i++ {
		x, y := random(0, screenWidth), random(0, screenHeight)
		d := dot{X: x, Y: y}
		a = append(a, d)
	}

	img = image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight))
	tmpImg, _ = ebiten.NewImage(screenWidth, screenHeight, ebiten.FilterNearest)

	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "Minimum Spanning Tree (Prim's Algorithm)"); err != nil {
		log.Fatal(err)
	}

}
