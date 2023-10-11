package atom

import (
	"ajz_xyz/experimental/computation/mfm-go"

	"github.com/nsf/termbox-go"
)

var ScissorsInfo = mfm.AtomInfo{
	Radius: 2,
	Rune:   'V',
	Fg:     termbox.ColorRed,
	USeq:   "✂",
	String: "Scissors",
}

func init() {
	ScissorsInfo.Func = Scissors
}

func Scissors(w *mfm.Window) {
	defer w.Stop()
	switch a := w.Atom(); a.Type {
	case nil:
		w.Swap(0)
	default:
		if w.Atom().Type.ID == PaperInfo.ID {
			w.Set(mfm.Atom{Type: &ScissorsInfo})
			return
		}
	}
}
