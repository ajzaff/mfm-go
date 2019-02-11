package main

import (
	"fmt"
	"time"

	"ajz_xyz/experimental/computation/mfm-go"
	"ajz_xyz/experimental/computation/mfm-go/atom"
	"ajz_xyz/experimental/computation/mfm-go/hook"
)

var cfg mfm.SimConfig

func init() {
	cfg.Seed = uint64(time.Now().UnixNano())
	cfg.Width, cfg.Height = 115, 35
}

func main() {
	sim := mfm.New(&cfg)
	stat := hook.NewStat(sim)

	sim.RegisterAtoms(
		&atom.DRegInfo,
		&atom.ResInfo,
		&atom.ForkInfo,
		&atom.SentryInfo,
		&atom.GermInfo,
		&atom.FernInfo,
		&atom.TreeInfo,
	)
	sim.RegisterHooks(stat)
	go func() {
		for {
			fmt.Printf("%c[2K\r", 27)
			fmt.Print(stat.String())
			time.Sleep(1 * time.Second)
		}
	}()
	sim.Set(mfm.C2D{X: 10, Y: 10}, mfm.Atom{Type: &atom.DRegInfo})
	sim.SetMode(0)
	sim.Run()
}
