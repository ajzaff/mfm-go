package main

import (
	"ajz_xyz/experimental/computation/mfm-go"
	"ajz_xyz/experimental/computation/mfm-go/hook"
)

func main() {
	sim := mfm.Sim{
		Bounds: &mfm.C2D{X: 115, Y: 36},
		Census: make(map[mfm.Atom]int),
		Sites:  make(map[mfm.C2D]mfm.Atom),
	}
	s := &hook.Stat{}
	sim.RegisterHooks(s, hook.NewTerm(s))
	sim.Run()
}
