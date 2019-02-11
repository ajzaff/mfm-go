package mfm

import (
	rand "ajz_xyz/numerics/random/xorshift64star"

	"github.com/nsf/termbox-go"

	"sync"
	"time"
)

// Sim represents a sparse atom simulation.
type Sim struct {
	Bounds *C2D
	Events int
	Census map[Atom]int
	Sites  map[C2D]Atom

	r  rand.Rand
	w  Window
	rw sync.RWMutex
	p  sync.Mutex
}

// Atom defines the interface for an atom in the simulation.
type Atom interface {
	Func(w *Window)
	String() string
	Rune() rune
	Color() termbox.Attribute
}

// Hook provides a extensible functionality the sim.
type Hook interface {
	Init(*Sim)
	Call(*Sim)
	Wait(*Sim)
}

// RegisterHooks registers and starts the hooks.
func (s *Sim) RegisterHooks(hooks ...Hook) {
	for _, h := range hooks {
		go func(h Hook) { // schedule hook
			h.Init(s)
			for {
				h.Wait(s)
				h.Call(s)
			}
		}(h)
	}
}

// Seed initializes the random number generator.
func (s *Sim) Seed(seed int64) { s.r.Seed(seed) }

// RLock acquires a reader exclusive lock.
func (s *Sim) RLock() { s.rw.RLock() }

// RUnlock releases a reader exclusive lock.
func (s *Sim) RUnlock() { s.rw.RUnlock() }

// Pause acquires the pause lock.
func (s *Sim) Pause() { s.p.Lock() }

// Unpause releases the pause lock.
func (s *Sim) Unpause() { s.p.Unlock() }

// Set updates the atom at the given position if valid.
func (s *Sim) Set(c C2D, a Atom) {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.set(c, a)
}

func (s *Sim) set(c C2D, a Atom) {
	if s.Bounds != nil && (c.X < 0 || c.Y < 0 || c.X >= s.Bounds.X || c.Y >= s.Bounds.Y) {
		return
	}
	if a, ok := s.Sites[c]; ok {
		s.Census[a]--
	}
	if a == nil {
		delete(s.Sites, c)
	} else {
		s.Sites[c] = a
		s.Census[a]++
	}
}

// Step advances the sim one time step.
func (s *Sim) Step() {
	s.step()
}

func (s *Sim) stepPausable() {
	s.p.Lock()
	defer s.p.Unlock()
	s.step()
}

func (s *Sim) step() {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.Events += len(s.Sites)
	for at, a := range s.Sites {
		s.w.loc = at
		a.Func(&s.w)
	}
}

// Clear resets the sim to an empty state.
func (s *Sim) Clear() {
	s.rw.Lock()
	defer s.rw.Unlock()
	s.Events = 0
	s.Census = make(map[Atom]int)
	s.Sites = make(map[C2D]Atom)
}

// Run begins the sim loop.
func (s *Sim) Run() {
	if s.r == 0 {
		s.r.Seed(time.Now().UnixNano())
	}
	s.p.Lock() // start out paused
	s.w.s = s
	for {
		s.stepPausable()
	}
}
