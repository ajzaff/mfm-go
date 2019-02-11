package atom

import (
	"github.com/nsf/termbox-go"

	"ajz_xyz/experimental/computation/mfm-go"
)

// FernInfo returns info about fern atoms.
var FernInfo = mfm.AtomInfo{
	Radius: 2,
	Rune:   '8',
	Fg:     termbox.ColorGreen,
	USeq:   "🌿",
	String: "Fern",
}

func init() {
	FernInfo.Func = Fern
}

// Fern implements snack atom behavior.
func Fern(w *mfm.Window) {
	if w.Atom().Type == nil && w.Roll(65535) { // grow
		w.Set(newGerm)
	}
	w.Stop()
}
