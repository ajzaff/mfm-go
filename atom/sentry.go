package atom

import (
	"github.com/nsf/termbox-go"

	"ajz_xyz/experimental/computation/mfm-go"
)

// SentryInfo contains Sentry atom info.
var SentryInfo = mfm.AtomInfo{
	Radius: 5,
	Rune:   'S',
	Fg:     termbox.ColorBlue,
	USeq:   "🛡",
	String: "Sentry",
}

func init() {
	SentryInfo.Func = Sentry
}

var winSize = uint64(41 * 3)

var (
	sentryAlert = Var{20, 10}
)

// Sentry defines the sentry atom behavior.
func Sentry(w *mfm.Window) {
	me := w.Self()
	a := w.Atom()
	switch {
	case a.Type == nil:
	case a.Type.ID == ResInfo.ID && w.Site() <= 8 && w.Roll(255):
		w.Set(mfm.Atom{Type: &SentryInfo, Value: 0})
	case a.Type.ID == ForkInfo.ID:
		sentryAlert.CountUint64(&me, winSize, winSize)
		w.Set(me)
		// w.SetValue(me.Value)
	case w.Site() == 0 && w.Roll(8191):
		w.Set(mfm.Atom{Type: &DRegInfo, Value: 0})
	}
	sentryAlert.CountUint64(&me, 0, winSize)
	if sentryAlert.ReadUint64(me) > 0 {
		w.Set(me)
	}
}
