package main

import (
	"flag"
	"sync"
	"time"

	"github.com/nsf/termbox-go"

	"ajz_xyz/experimental/computation/mfm-go"
	"ajz_xyz/experimental/computation/mfm-go/atom"
	"ajz_xyz/experimental/computation/mfm-go/hook"
)

var (
	enableUnicodeAtoms = flag.Bool("enable_unicode", false,
		"whether to enable unicode atom output")
)

var (
	cfg       mfm.SimConfig
	eventChan chan termbox.Event
)

func init() {
	flag.Parse()

	var wg sync.WaitGroup
	wg.Add(1)
	cfg.Mode = mfm.ModePause
	cfg.Seed = uint64(time.Now().UnixNano())
	go func() {
		// Initialize term hook and termbox library and poll.
		// We need to do everything in one goroutine with termbox.
		if err := termbox.Init(); err != nil {
			panic(err)
		}
		cfg.Width, cfg.Height = termbox.Size()
		wg.Done()

		eventChan = make(chan termbox.Event)
		for {
			eventChan <- termbox.PollEvent()
		}
	}()
	wg.Wait()
}

func main() {
	sim := mfm.New(&cfg)
	stat := hook.NewStat(sim)
	term := hook.NewTerm(&hook.TermConfig{
		Stat:      stat,
		Sim:       sim,
		EventChan: eventChan,

		EnableUnicodeAtoms: *enableUnicodeAtoms,
	})

	sim.RegisterAtoms(
		&atom.DRegInfo,
		&atom.ResInfo,
		&atom.ForkInfo,
		&atom.SentryInfo,
		&atom.GermInfo,
		&atom.FernInfo,
		&atom.TreeInfo,
		&atom.RockInfo,
		&atom.PaperInfo,
		&atom.ScissorsInfo,
	)
	sim.RegisterHooks(stat, term)
	sim.Run()
}
