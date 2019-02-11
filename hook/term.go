package hook

import (
	"os"
	"time"

	"github.com/nsf/termbox-go"

	"ajz_xyz/experimental/computation/mfm-go"
	"ajz_xyz/experimental/computation/mfm-go/atom"
)

// Term is a hook providing terminal based events.
type Term struct {
	s      *Stat
	sim    *mfm.Sim
	cur    struct{ x, y int }
	v1, v2 struct{ x, y int }
	mode   Mode
	ev     chan termbox.Event

	EnableUnicodeAtoms bool
}

// TermConfig provides condiguration to create a new terminal hook.
type TermConfig struct {
	Stat *Stat
	Sim  *mfm.Sim

	EnableUnicodeAtoms bool
	EventChan          chan termbox.Event
}

// NewTerm creates a new terminal hook.
// Termbox must be initialized before calling NewTerm.
func NewTerm(c *TermConfig) *Term {
	t := &Term{
		s:   c.Stat,
		sim: c.Sim,
		ev:  c.EventChan,

		EnableUnicodeAtoms: c.EnableUnicodeAtoms,
	}

	return t
}

// Mode enumerates the terminal mode flags.
type Mode int

// Enumerate mode flags.
const (
	ModeVisual = 1 << iota
	ModeReplace
)

func (t *Term) setMode(m Mode) {
	v := t.mode
	t.mode = m
	if (v^t.mode)&ModeVisual != 0 && t.mode&ModeVisual != 0 {
		// ModeVisual: enable visual mode.
		t.v1.x, t.v1.y = t.cur.x, t.cur.y
		t.v2.x, t.v2.y = t.cur.x, t.cur.y
	}
}

// Wait is called when the hook is in the waiting state.
func (*Term) Wait() { time.Sleep(10 * time.Millisecond) }

func (t *Term) updateCursor(dx, dy int) {
	t.mode |= ModeReplace
	t.cur.x += dx
	t.cur.y += dy
	if t.cur.x < 0 {
		t.cur.x = 0
	}
	if t.cur.y < 0 {
		t.cur.y = 0
	}
	w, h := termbox.Size()
	if t.cur.x >= w {
		t.cur.x = w - 1
	}
	if t.cur.y >= h {
		t.cur.y = h - 1
	}
	if t.mode&ModeVisual != 0 {
		t.v2.x, t.v2.y = t.cur.x, t.cur.y
	}
}

func (t *Term) getVisualRect() (x0, y0, x1, y1 int) {
	x0, y0, x1, y1 = t.v1.x, t.v1.y, t.v2.x, t.v2.y
	if x1 < x0 {
		x0, x1 = x1, x0
	}
	if y1 < y0 {
		y0, y1 = y1, y0
	}
	return
}

func (t *Term) isVisualArea(x, y int) bool {
	if t.mode&ModeVisual != 0 {
		x0, y0, x1, y1 := t.getVisualRect()
		return x0 <= x && x <= x1 && y0 <= y && y <= y1
	}
	return false
}

// Call outputs stat information to the terminal.
func (t *Term) draw() {
	w, h := termbox.Size()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	state := t.sim.State()
	for c, a := range state.Sites {
		x, y := c.X, c.Y
		bg := termbox.ColorDefault
		if t.isVisualArea(x, y) {
			bg = termbox.ColorWhite
		}
		if t.EnableUnicodeAtoms {
			for i, c := range a.Type.USeq {
				termbox.SetCell(x+i, y, c, termbox.ColorDefault, bg)
			}
		} else {
			termbox.SetCell(x, y, a.Type.Rune, a.Type.Fg, bg)
		}
	}

	// Write atom population stats.
	stats := t.s.String()
	for i, r := range t.s.String() {
		termbox.SetCell(w-len(stats)+i, h-1, r,
			termbox.ColorDefault, termbox.ColorDefault)
	}

	var mode string

	switch {
	case t.mode&ModeVisual != 0:
		mode = "-- VISUAL --"
	case t.mode&ModeReplace != 0:
		mode = "-- REPLACE --"
	case state.Mode&mfm.ModePause != 0:
		mode = "-- PAUSED --"
	}
	for i, r := range mode {
		termbox.SetCell(i, h-1, r,
			termbox.ColorDefault|termbox.AttrBold, termbox.ColorDefault)
	}

	if t.mode&(ModeVisual|ModeReplace) != 0 {
		termbox.SetCursor(t.cur.x, t.cur.y)
	} else {
		termbox.HideCursor()
	}

	termbox.Flush()
}

func (t *Term) handleVisualPaste(a mfm.Atom) {
	if t.mode&ModeVisual != 0 {
		x0, y0, x1, y1 := t.getVisualRect()
		for x := x0; x <= x1; x++ {
			for y := y0; y <= y1; y++ {
				t.sim.Set(mfm.C2D{x, y}, a)
			}
		}
		t.setMode(t.mode &^ ModeVisual)
	} else if t.mode&ModeReplace != 0 {
		t.sim.Set(mfm.C2D{t.cur.x, t.cur.y}, a)
	}
}

// Call outputs stat information to the terminal.
func (t *Term) Call() {
	select {
	case e := <-t.ev:
		if e.Type == termbox.EventKey {
			switch {
			case e.Key == termbox.KeyCtrlC:
				termbox.Close()
				os.Exit(0)
			case e.Key == termbox.KeySpace:
				t.sim.SetMode(t.sim.Mode() ^ mfm.ModePause)
			case e.Ch == '.':
				t.sim.SetMode(t.sim.Mode() | mfm.ModePause)
				t.sim.Step()
			case e.Key == termbox.KeyCtrlX: // clear
				t.sim.Reset()
			case e.Key == termbox.KeyEsc:
				t.setMode(t.mode & ^(ModeVisual | ModeReplace))
			case e.Ch == 'r':
				if t.mode&(ModeReplace|ModeVisual) == 0 {
					t.setMode(t.mode | ModeReplace)
				} else {
					t.handleVisualPaste(mfm.Atom{&atom.ResInfo, 0})
				}
			case e.Ch == 'v':
				t.setMode(t.mode ^ ModeVisual)
			case e.Ch == 'H':
				t.updateCursor(-10, 0)
			case e.Ch == 'K':
				t.updateCursor(10, 0)
			case e.Ch == 'U':
				t.updateCursor(0, -10)
			case e.Ch == 'M':
				t.updateCursor(0, 10)
			case e.Ch == 'h' || e.Key == termbox.KeyArrowLeft:
				t.updateCursor(-1, 0)
			case e.Ch == 'k' || e.Key == termbox.KeyArrowRight:
				t.updateCursor(1, 0)
			case e.Ch == 'u' || e.Key == termbox.KeyArrowUp:
				t.updateCursor(0, -1)
			case e.Ch == 'm' || e.Key == termbox.KeyArrowDown:
				t.updateCursor(0, 1)
			case e.Ch == 'D':
				t.handleVisualPaste(mfm.Atom{&atom.DRegInfo, 0})
			case e.Ch == 'x':
				t.handleVisualPaste(mfm.Atom{})
			case e.Ch == 'X':
				t.handleVisualPaste(mfm.Atom{&atom.ForkInfo, 0})
			case e.Ch == 'S':
				t.handleVisualPaste(mfm.Atom{&atom.SentryInfo, 0})
			case e.Ch == '7':
				t.handleVisualPaste(mfm.Atom{&atom.GermInfo, 0})
			case e.Ch == '8':
				t.handleVisualPaste(mfm.Atom{&atom.FernInfo, 0})
			case e.Ch == '9':
				t.handleVisualPaste(mfm.Atom{&atom.TreeInfo, 0})
			case e.Ch == 'R':
				t.handleVisualPaste(mfm.Atom{&atom.RabbitInfo, 0})
			case e.Ch == 'E':
				t.handleVisualPaste(mfm.Atom{&atom.HerbasaurInfo, 0})
			case e.Ch == '^':
				t.v2.x, t.v2.y, t.cur.x, t.cur.y = 0, 0, 0, 0
			case e.Ch == '$':
				w, h := termbox.Size()
				t.v2.x, t.v2.y, t.cur.x, t.cur.y = w-1, h-1, w-1, h-1
			}
		}
	default:
		t.draw()
	}
}
