package hook

import (
	"ajz_xyz/experimental/computation/mfm-go"
	"ajz_xyz/experimental/computation/mfm-go/atom"

	"github.com/nsf/termbox-go"

	"bytes"
	"fmt"
	"os"
	"time"
)

// Term is a hook providing terminal based events.
type Term struct {
	s      *Stat
	cur    mfm.C2D
	v1, v2 mfm.C2D
	mode   Mode
	ev     chan termbox.Event
}

// Mode enumerates the terminal mode flags.
type Mode int

// Enumerate mode flags.
const (
	ModePause Mode = 1 << iota
	ModeDefault
	ModeVisual
	ModeReplace
)

func (t *Term) setMode(s *mfm.Sim, m Mode) {
	v := t.mode
	t.mode = m
	if (v^t.mode)&ModePause != 0 { // ModePause change
		if t.mode&ModePause != 0 {
			s.Pause()
		} else {
			s.Unpause()
		}
	}
	if (v^t.mode)&ModeVisual != 0 && t.mode&ModeVisual != 0 {
		// ModeVisual: enable visual mode.
		t.v1.SetC2D(t.cur)
		t.v2.SetC2D(t.cur)
	}
}

// NewTerm creates a new terminal.
func NewTerm(s *Stat) *Term {
	return &Term{s: s, mode: ModePause}
}

// Init is called when the hook is registered.
func (t *Term) Init(s *mfm.Sim) {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	w, h := termbox.Size()
	s.Bounds.Set(w, h) // set bounds.
	go func() {
		t.ev = make(chan termbox.Event)
		for {
			t.ev <- termbox.PollEvent()
		}
	}()
}

// Wait is called when the hook is in the waiting state.
func (*Term) Wait(*mfm.Sim) {
	time.Sleep(10 * time.Millisecond)
}

func (t *Term) updateCursor(dx, dy int) {
	t.mode |= ModeReplace
	t.cur.X += dx
	t.cur.Y += dy
	if t.cur.X < 0 {
		t.cur.X = 0
	}
	if t.cur.Y < 0 {
		t.cur.Y = 0
	}
	w, h := termbox.Size()
	if t.cur.X >= w {
		t.cur.X = w - 1
	}
	if t.cur.Y >= h {
		t.cur.Y = h - 1
	}
	if t.mode&ModeVisual != 0 {
		t.v2.SetC2D(t.cur)
	}
}

func (t *Term) isVisualArea(x, y int) bool {
	if t.mode&ModeVisual != 0 {
		x0, y0, x1, y1 := t.v1.X, t.v1.Y, t.v2.X, t.v2.Y
		if x1 < x0 {
			x0, x1 = x1, x0
		}
		if y1 < y0 {
			y0, y1 = y1, y0
		}
		return x0 <= x && x <= x1 && y0 <= y && y <= y1
	}
	return false
}

// Call outputs stat information to the terminal.
func (t *Term) draw(s *mfm.Sim) {
	var c mfm.C2D
	w, h := termbox.Size()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	s.RLock()
	defer s.RUnlock()

	for x := 0; x < w; x++ {
		c.X = x
		for y := 0; y < h; y++ {
			c.Y = y
			bg := termbox.ColorDefault
			if t.isVisualArea(x, y) {
				bg = termbox.ColorWhite
			}
			if a, ok := s.Sites[c]; ok {
				termbox.SetCell(x, y, a.Rune(), a.Color(), bg)
			} else {
				termbox.SetCell(x, y, ' ', termbox.ColorDefault, bg)
			}
		}
	}

	// Write atom population stats.
	buf := bytes.NewBufferString("")
	buf.WriteString(fmt.Sprintf("%d", t.s.Census[" "]))
	for _, n := range t.s.Names {
		buf.WriteString(fmt.Sprintf(" %s%d", n, t.s.Census[n]))
	}
	buf.WriteString(fmt.Sprintf(" %d%% %d(%dK/s)",
		t.s.Fullness, s.Events, t.s.KEPS))

	for i, r := range buf.String() {
		termbox.SetCell(w-buf.Len()+i, h-1, r,
			termbox.ColorDefault, termbox.ColorDefault)
	}

	var mode string

	switch {
	case t.mode&ModeVisual != 0:
		mode = "-- VISUAL --"
	case t.mode&ModeReplace != 0:
		mode = "-- REPLACE --"
	case t.mode&ModePause != 0:
		mode = "-- PAUSED --"
	}
	for i, r := range mode {
		termbox.SetCell(i, h-1, r,
			termbox.ColorDefault|termbox.AttrBold, termbox.ColorDefault)
	}

	if t.mode&(ModeVisual|ModeReplace) != 0 {
		termbox.SetCursor(t.cur.X, t.cur.Y)
	} else {
		termbox.HideCursor()
	}

	termbox.Flush()
}

// Call outputs stat information to the terminal.
func (t *Term) Call(s *mfm.Sim) {
	select {
	case e := <-t.ev:
		if e.Type == termbox.EventKey {
			switch {
			case e.Key == termbox.KeyCtrlC:
				termbox.Close()
				os.Exit(0)
			case e.Key == termbox.KeySpace:
				t.setMode(s, t.mode^ModePause)
			case e.Ch == '.':
				t.setMode(s, t.mode|ModePause)
				s.Step()
			case e.Key == termbox.KeyCtrlX:
				t.setMode(s, t.mode|ModePause)
				s.Clear()
			case e.Key == termbox.KeyEsc:
				t.setMode(s, t.mode & ^(ModeVisual|ModeReplace))
			case e.Ch == 'r':
				if t.mode&ModeVisual != 0 {
					x0, y0, x1, y1 := t.v1.X, t.v1.Y, t.v2.X, t.v2.Y
					if x1 < x0 {
						x0, x1 = x1, x0
					}
					if y1 < y0 {
						y0, y1 = y1, y0
					}
					var c mfm.C2D
					for x := x0; x <= x1; x++ {
						c.X = x
						for y := y0; y <= y1; y++ {
							c.Y = y
							s.Set(c, atom.Res(0))
						}
					}
					t.setMode(s, t.mode&^ModeVisual)
				} else if t.mode&ModeReplace != 0 {
					s.Set(t.cur, atom.Res(0))
				} else {
					t.setMode(s, t.mode^ModeReplace)
				}
			case e.Ch == 'v':
				t.setMode(s, t.mode^ModeVisual)
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
				if t.mode&ModeVisual != 0 {
					x0, y0, x1, y1 := t.v1.X, t.v1.Y, t.v2.X, t.v2.Y
					if x1 < x0 {
						x0, x1 = x1, x0
					}
					if y1 < y0 {
						y0, y1 = y1, y0
					}
					var c mfm.C2D
					for x := x0; x <= x1; x++ {
						c.X = x
						for y := y0; y <= y1; y++ {
							c.Y = y
							s.Set(c, atom.DReg(0))
						}
					}
					t.setMode(s, t.mode&^ModeVisual)
				} else if t.mode&ModeReplace != 0 {
					s.Set(t.cur, atom.DReg(0))
				}
			case e.Ch == 'x':
				if t.mode&ModeVisual != 0 {
					// Clear the visual region.
					x0, y0, x1, y1 := t.v1.X, t.v1.Y, t.v2.X, t.v2.Y
					if x1 < x0 {
						x0, x1 = x1, x0
					}
					if y1 < y0 {
						y0, y1 = y1, y0
					}
					var c mfm.C2D
					for x := x0; x <= x1; x++ {
						c.X = x
						for y := y0; y <= y1; y++ {
							c.Y = y
							s.Set(c, nil)
						}
					}
					t.setMode(s, t.mode&^ModeVisual)
				} else if t.mode&ModeReplace != 0 {
					s.Set(t.cur, nil)
				}
			case e.Ch == 'X':
				if t.mode&ModeVisual != 0 {
					x0, y0, x1, y1 := t.v1.X, t.v1.Y, t.v2.X, t.v2.Y
					if x1 < x0 {
						x0, x1 = x1, x0
					}
					if y1 < y0 {
						y0, y1 = y1, y0
					}
					var c mfm.C2D
					for x := x0; x <= x1; x++ {
						c.X = x
						for y := y0; y <= y1; y++ {
							c.Y = y
							s.Set(c, atom.Fork(0))
						}
					}
					t.setMode(s, t.mode&^ModeVisual)
				} else if t.mode&ModeReplace != 0 {
					s.Set(t.cur, atom.Fork(0))
				}
			case e.Ch == '^':
				t.setMode(s, t.mode^ModeVisual)
				t.v2.Set(0, 0)
				t.cur.Set(0, 0)
			case e.Ch == '$':
				t.setMode(s, t.mode^ModeVisual)
				w, h := termbox.Size()
				t.v2.Set(w-1, h-1)
				t.cur.Set(w-1, h-1)
			}
		}
	default:
		t.draw(s)
	}
}
