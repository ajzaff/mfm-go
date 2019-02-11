package bench

import (
	"testing"
	"time"

	"ajz_xyz/experimental/computation/mfm-go"
	"ajz_xyz/experimental/computation/mfm-go/atom"
)

func newSim() *mfm.Sim {
	sim := mfm.New(&mfm.SimConfig{
		Width:  10000,
		Height: 10000,
	})
	sim.RegisterAtoms(
		&atom.DRegInfo,
		&atom.ResInfo,
		&atom.SentryInfo,
	)
	return sim
}

// BenchmarkDreg runs a DReg atom benchmark.
func BenchmarkDReg(b *testing.B) {
	sim := newSim()
	sim.Set(mfm.C2D{39, 9}, mfm.Atom{Type: &atom.DRegInfo})
	go sim.Run()
	<-time.After(600 * time.Second)
}

// BenchmarkSentry runs a sentry atom benchmark.
func BenchmarkSentry(b *testing.B) {
	sim := newSim()
	sim.Set(mfm.C2D{39, 9}, mfm.Atom{Type: &atom.DRegInfo})
	for i := 0; i < 100; i++ {
		sim.Set(mfm.C2D{i, 39}, mfm.Atom{Type: &atom.SentryInfo})
	}
	go sim.Run()
	<-time.After(600 * time.Second)
}
