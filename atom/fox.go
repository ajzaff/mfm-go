package atom

import (
	"github.com/nsf/termbox-go"

	"ajz_xyz/experimental/computation/mfm-go"
)

//FoxInfo defines the info for a fox atom.
var FoxInfo = mfm.AtomInfo{
	Radius: 2,
	Rune:   'F',
	Fg:     termbox.ColorYellow,
	USeq:   "F",
	String: "Fox",
}

func init() {
	FoxInfo.Func = Fox
}

// Fox implements a fox atom.
func Fox(w *mfm.Window) {}
