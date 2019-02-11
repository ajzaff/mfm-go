package atom

import (
	"github.com/nsf/termbox-go"

	"ajz_xyz/experimental/computation/mfm-go"
)

// ForkInfo defines the info for a forking atom.
var ForkInfo = mfm.AtomInfo{
	Radius: 5,
	Rune:   'X',
	Fg:     termbox.ColorRed,
	USeq:   "💥",
	String: "Fork",
}

func init() {
	ForkInfo.Func = Fork
}

var newFork = mfm.Atom{Type: &ForkInfo, Value: 0}

// Fork defines the fork function.
func Fork(w *mfm.Window) {
	w.Set(newFork)
}
