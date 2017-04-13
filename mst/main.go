package main

import (
	"math/rand"
	"time"

	"github.com/crgimenes/graphos/coreScreen"
)

var (
	cs *coreScreen.Instance
)

func update(screen *coreScreen.Instance) error {
	return nil
}

func main() {
	rand.Seed(time.Now().Unix())

	cs = coreScreen.Get()
	cs.Update = update
	cs.Title = "Minimum Spanning Tree (Prim's Algorithm)"
	cs.CurrentColor = 0x9

	cs.Run()
}
