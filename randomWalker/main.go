package main

import (
	"math/rand"
	"time"

	"github.com/crgimenes/graphos/coreGame"
)

type walker struct {
	X       int
	Y       int
	lastDir int
}

var (
	cg      *coreGame.Instance
	walkers []walker
)

var color byte

func getNextColor() byte {
	color++
	if color > 15 {
		color = 1
	}
	return color
}

func update(screen *coreGame.Instance) error {

	for i := 0; i < len(walkers); i++ {
		x := random(1, 5)
		switch x {
		case 1:
			walkers[i].X++
		case 2:
			walkers[i].X--
		case 3:
			walkers[i].Y++
		case 4:
			walkers[i].Y--
		}
	}

	// --=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-=-

	for i := 0; i < len(walkers); i++ {
		screen.CurrentColor = getNextColor()
		//screen.DrawFilledCircle(walkers[i].X, walkers[i].Y, 1)
		screen.DrawPix(walkers[i].X, walkers[i].Y)
	}
	return nil
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func main() {
	rand.Seed(time.Now().Unix())

	cg = coreGame.Get()
	cg.Width = 800
	cg.Height = 600
	cg.Scale = 1
	cg.ScreenHandler = update
	cg.Title = "Random Walker"
	cg.CurrentColor = 0x0

	for i := 0; i < 60; i++ {
		w := walker{
			X: random(10, cg.Width-10),
			Y: random(10, cg.Height-10),
		}
		walkers = append(walkers, w)
	}

	cg.Run()
}
