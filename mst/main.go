package main

import (
	"math/rand"
	"time"

	"github.com/crgimenes/graphos/coreGame"
)

var (
	cs *coreScreen.Instance
)

func update(screen *coreScreen.Instance) error {
	return nil
}

func main() {
	rand.Seed(time.Now().Unix())

	cs = coreGame.Get()
	cs.ScreenHandler = update
	cs.Title = "Minimum Spanning Tree (Prim's Algorithm)"
	cs.CurrentColor = 0x9

	cs.Run()
}
