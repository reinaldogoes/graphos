package graphos

import (
	"image"
	"log"

	"github.com/crgimenes/graphos/fonts"
	"github.com/hajimehoshi/ebiten"
)

type Instance struct {
	Border        int
	Height        int
	Width         int
	Scale         float64
	CurrentColor  byte
	uTime         int
	updateScreen  bool
	tmpScreen     *ebiten.Image
	img           *image.RGBA
	ScreenHandler func(*Instance) error
	Title         string
	Font          struct {
		Height int
		Width  int
		Bitmap [][]byte
	}
}

var cg *Instance = nil // Current Instance

func Get() *Instance {
	if cg == nil {
		cg = &Instance{
			Scale:  2,
			Border: 0,
			Height: 240,
			Width:  320,
		}
	}
	return cg
}

var Colors = []struct {
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

func MergeColorCode(b, f byte) byte {
	return (f & 0xff) | (b << 4)
}

func update(screen *ebiten.Image) error {

	if cg.ScreenHandler != nil {
		err := cg.ScreenHandler(cg)
		if err != nil {
			return err
		}
	}

	if cg.updateScreen {
		cg.tmpScreen.ReplacePixels(cg.img.Pix)
		cg.updateScreen = false
	}

	screen.DrawImage(cg.tmpScreen, nil)
	cg.uTime++
	return nil
}

func (p *Instance) Run() {

	{
		var f fonts.Expert118x8
		f.Load()
		p.Font.Bitmap = f.Bitmap
		p.Font.Height = f.Height
		p.Font.Width = f.Width
	}

	p.img = image.NewRGBA(image.Rect(0, 0, p.Width, p.Height))
	p.tmpScreen, _ = ebiten.NewImage(p.Width, p.Height, ebiten.FilterNearest)

	p.Clear()
	p.updateScreen = true

	if err := ebiten.Run(update, p.Width, p.Height, p.Scale, p.Title); err != nil {
		log.Fatal(err)
	}
}

func (p *Instance) DrawPix(x, y int) {
	x += p.Border
	y += p.Border
	if x < p.Border || y < p.Border || x >= p.Width-p.Border || y >= p.Height-p.Border {
		return
	}
	pos := 4*y*p.Width + 4*x
	p.img.Pix[pos] = Colors[p.CurrentColor].R
	p.img.Pix[pos+1] = Colors[p.CurrentColor].G
	p.img.Pix[pos+2] = Colors[p.CurrentColor].B
	p.img.Pix[pos+3] = 0xff
	p.updateScreen = true
}

func (p *Instance) Clear() {
	for i := 0; i < p.Height*p.Width*4; i += 4 {
		p.img.Pix[i] = Colors[p.CurrentColor].R
		p.img.Pix[i+1] = Colors[p.CurrentColor].G
		p.img.Pix[i+2] = Colors[p.CurrentColor].B
		p.img.Pix[i+3] = 0xff
	}
}
