package atom

import (
	"ajz_xyz/experimental/computation/mfm-go"

	"github.com/nsf/termbox-go"
)

// Fork is a forking atom type.
type Fork int

// Func defines the forking atom behavior.
func (a Fork) Func(w *mfm.Window) {
	var abuf [41]mfm.Atom
	for i := range mfm.InitialWindow {
		abuf[i] = Fork(0)
	}
	w.Write(mfm.InitialWindow[:], abuf[:])
}

// Rune returns the rune for this atom.
func (a Fork) Rune() rune { return 'X' }

// Color returns the terminal foreground attribute.
func (a Fork) Color() termbox.Attribute { return termbox.ColorRed }

func (a Fork) String() string { return "fork" }
