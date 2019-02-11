package atom

import (
	"ajz_xyz/experimental/computation/mfm-go"

	"github.com/nsf/termbox-go"
)

// Res is a resource atom type.
type Res int

// Func implements resource atom behavior.
func (a Res) Func(w *mfm.Window) {
	if s, ok := w.Diffuse(); ok {
		w.Write([]mfm.Site{mfm.Me, s},
			[]mfm.Atom{nil, a})
	}
}

// Rune returns the rune for this atom.
func (a Res) Rune() rune { return 'r' }

// Color returns the terminal foreground attribute.
func (a Res) Color() termbox.Attribute { return termbox.ColorYellow }

func (a Res) String() string { return "res" }
