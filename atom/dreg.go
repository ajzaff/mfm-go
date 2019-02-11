package atom

import (
	"ajz_xyz/experimental/computation/mfm-go"

	"github.com/nsf/termbox-go"
)

// DReg atom implements a dynamic space regulator atom.
// Roles include producing Res, random destruction,
// and occasional reproduction.
type DReg int

// Func implements the DReg atom behavior.
func (a DReg) Func(w *mfm.Window) {
	n := 0
	var abuf [10]mfm.Atom
	var sbuf [10]mfm.Site
	err := w.At(mfm.NeighborDeltas, func(s mfm.Site, a mfm.Atom) error {
		if w.Valid(s) {
			if a == nil {
				if w.Roll(511) {
					// Create new Res.
					sbuf[n] = s
					abuf[n] = Res(0)
					n++
				} else if w.Roll(2047) {
					// Reproduce new DReg.
					sbuf[n] = s
					abuf[n] = DReg(0)
					n++
				}
			} else if _, ok := a.(DReg); (ok && w.Roll(7)) || w.Roll(255) {
				// Destroy an adjacent DReg or general destruction.
				sbuf[n] = s
				abuf[n] = nil
				n++
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	if s, ok := w.Diffuse(); ok {
		sbuf[n], abuf[n] = mfm.Me, nil
		n++
		sbuf[n], abuf[n] = s, a
		n++
	}
	if n > 0 {
		w.Write(sbuf[:n], abuf[:n])
	}
}

// Rune returns the rune for this atom.
func (a DReg) Rune() rune { return 'D' }

// Color returns the terminal foreground attribute.
func (a DReg) Color() termbox.Attribute { return termbox.ColorMagenta }

func (a DReg) String() string { return "dreg" }
