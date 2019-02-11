package hook

import (
	"ajz_xyz/experimental/computation/mfm-go"
	"sort"

	"time"
)

// Stat simulation stats
type Stat struct {
	KEPS     int
	Fullness int
	Census   map[string]int
	Names    []string

	lastEvents int
	lastClock  int64
}

// Init is called when the hook is registered.
func (h *Stat) Init(*mfm.Sim) {
	h.Census = make(map[string]int)
	h.Names = make([]string, 0)
}

// Wait is called when the hook is in the waiting state.
func (Stat) Wait(*mfm.Sim) {
	time.Sleep(1 * time.Second)
}

// Call outputs stat information to the terminal.
func (h *Stat) Call(s *mfm.Sim) {
	now := time.Now().UnixNano()

	s.RLock()
	defer s.RUnlock()

	h.Names = h.Names[:0]
	for n := range h.Census {
		delete(h.Census, n)
	}
	for a, v := range s.Census {
		h.Census[string(a.Rune())] = v
		h.Names = append(h.Names, string(a.Rune()))
	}
	sort.Strings(h.Names)

	// Write bounds dependent stats.
	h.Census[" "] = (s.Bounds.X * s.Bounds.Y) - len(s.Sites)
	h.Fullness = 100 * len(s.Sites) / (s.Bounds.X * s.Bounds.Y)

	// Compute perf stats.
	events := 0
	if h.lastEvents <= s.Events {
		events = s.Events - h.lastEvents
		second := float64(now-h.lastClock) / float64(time.Second)
		h.KEPS = int(float64(events) / (1000 * second))
	}
	h.lastEvents = s.Events
	h.lastClock = now
}
