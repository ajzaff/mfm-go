package pprof

import (
	"testing"
	"time"

	"ajz_xyz/experimental/computation/mfm-go"
	"ajz_xyz/experimental/computation/mfm-go/atom"
)

// BenchmarkDreg runs a dreg benchmark for 10 seconds.
func BenchmarkDreg(b *testing.B) {
	sim := mfm.Sim{
		Bounds: &mfm.C2D{X: 8000, Y: 3000},
		Census: make(map[mfm.Atom]int),
		Sites:  make(map[mfm.C2D]mfm.Atom),
	}
	sim.Set(mfm.C2D{X: 39, Y: 9}, atom.DReg(0))
	go sim.Run()
	<-time.After(60 * time.Second)
}
