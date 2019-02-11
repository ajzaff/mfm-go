package mfm

import "sync"

// SimConfig holds configuration used to create a new simulation.
type SimConfig struct {
	Mode   SimMode
	Seed   uint64
	Width  int
	Height int
}

// SimState contains public data related to the running the simulation.
type SimState struct {
	Sites  map[C2D]Atom
	Mode   SimMode
	Events int
	Width  int
	Height int
}

// SimMode defines various running modes of the sim.
type SimMode int

// Enumerate sim states.
const (
	ModeRun           = 0
	ModePause SimMode = 1 << iota
	ModeYield
)

// Sim represents a sparse atom simulation.
type Sim struct {
	state  SimState
	win    Window
	nextID uint64
	scond  sync.Cond
	smutex sync.Mutex
}

// New creates a new mfm sim.
func New(c *SimConfig) *Sim {
	sim := &Sim{state: SimState{
		Mode:   c.Mode,
		Width:  c.Width,
		Height: c.Height,
		Sites:  make(map[C2D]Atom),
	}}
	sim.scond.L = &sim.smutex
	sim.win.rand.Seed(c.Seed)
	return sim
}

func (s *Sim) SetMode(t SimMode) {
	s.scond.L.Lock()
	defer s.scond.L.Unlock()
	s.state.Mode = t
	if s.state.Mode == 0 {
		s.scond.Signal()
	}
}
func (s *Sim) Mode() SimMode { return s.state.Mode }
func (s *Sim) State() SimState {
	state := SimState{
		Sites:  make(map[C2D]Atom),
		Mode:   s.state.Mode,
		Events: s.state.Events,
		Width:  s.state.Width,
		Height: s.state.Height,
	}
	for c, a := range s.state.Sites { // copy sites
		state.Sites[c] = a
	}
	return state
}
func (s *Sim) Reset() { s.state.Sites = make(map[C2D]Atom) }
func (s *Sim) valid(c C2D) bool {
	return (s.state.Width == 0 && s.state.Height == 0) ||
		(c.X >= 0 && c.Y >= 0 && c.X < s.state.Width && c.Y < s.state.Height)
}
func (s *Sim) Set(c C2D, a Atom) (ok bool) {
	if s.valid(c) {
		ok = true
		if a.Type == nil {
			delete(s.state.Sites, c)
		} else {
			s.state.Sites[c] = a
		}
	}
	return
}
func (s *Sim) Get(c C2D) (a Atom, ok bool) {
	return s.state.Sites[c], s.valid(c)
}
func (s *Sim) Step() { // must be yielded or paused
	var c2 C2D

	s.state.Events += len(s.state.Sites)
	for c, a := range s.state.Sites {
		n := r[a.Type.Radius]
		if n == 0 {
			continue
		}
		s.win.loc = c
		s.win.me = a
		site := Site(uint32((uint64(uint32(s.win.rand.Uint64())) * uint64(n)) >> 32))
		c2.Set(site)
		c2.Add(c2, s.win.loc)
		if !s.valid(c2) {
			continue
		}
		s.win.Reset()
		s.win.s = site
		s.win.a, _ = s.state.Sites[c2]
		a.Type.Func(&s.win)
		var dst C2D
		dst.Set(s.win.mut.dst)
		dst.Add(dst, s.win.loc)
		switch {
		case s.win.mut.mode&mutSet != 0:
			s.Set(c2, s.win.mut.atom)
		case s.win.mut.mode&mutMove != 0:
			a, _ := s.Get(c2)
			s.Set(dst, a)
		case s.win.mut.mode&mutSwap != 0:
			a, _ := s.Get(c2)
			if a2, ok := s.Get(dst); ok {
				s.Set(dst, a)
				s.Set(c2, a2)
			}
		}
	}
}
func (s *Sim) Run() {
	for {
		s.scond.L.Lock()
		for s.state.Mode != 0 {
			s.scond.Wait()
		}
		s.Step()
		s.scond.L.Unlock()
	}
}
func (s *Sim) RegisterAtoms(atoms ...*AtomInfo) {
	for _, a := range atoms {
		s.nextID++
		a.ID = s.nextID
	}
}
func (s *Sim) RegisterHooks(hooks ...Hook) {
	var hookMu sync.Mutex
	for _, h := range hooks {
		go func(h Hook) { // schedule hook
			for {
				h.Wait()
				hookMu.Lock()
				s.SetMode(s.state.Mode ^ ModeYield)
				h.Call()
				s.SetMode(s.state.Mode ^ ModeYield)
				hookMu.Unlock()
			}
		}(h)
	}
}

// Hook provides a extensible functionality the simulation.
type Hook interface {
	Call()
	Wait()
}
