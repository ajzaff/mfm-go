package atom

import (
	"ajz_xyz/experimental/computation/mfm-go"

	"github.com/nsf/termbox-go"
)

var RockInfo = mfm.AtomInfo{
	Radius: 2,
	Rune:   'o',
	Fg:     termbox.ColorBlue,
	USeq:   "🪨",
	String: "Rock",
}

func init() {
	RockInfo.Func = Rock
}

func Rock(w *mfm.Window) {
	defer w.Stop()
	switch a := w.Atom(); a.Type {
	case nil:
		w.Swap(0)
	default:
		if w.Atom().Type.ID == ScissorsInfo.ID {
			w.Set(mfm.Atom{Type: &RockInfo})
			return
		}
	}
}
