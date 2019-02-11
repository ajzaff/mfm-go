package hook

import (
	"bytes"
	"fmt"
	"sort"
	"sync"
	"time"

	"ajz_xyz/experimental/computation/mfm-go"
)

// Stat simulation stats
type Stat struct {
	s          *mfm.Sim
	keps       int
	fullness   int
	census     map[string]int
	names      []string
	lastEvents int
	lastClock  int64
	m          sync.RWMutex
}

// NewStat creates a new stat hook.
func NewStat(s *mfm.Sim) *Stat {
	return &Stat{
		s:      s,
		census: make(map[string]int),
		names:  make([]string, 0),
	}
}

// Wait is called when the hook is in the waiting state.
func (*Stat) Wait() { time.Sleep(1 * time.Second) }

func (h *Stat) String() string {
	h.m.RLock()
	defer h.m.RUnlock()

	buf := bytes.NewBufferString("")
	buf.WriteString(fmt.Sprintf("%d", h.census[" "]))
	for _, n := range h.names {
		buf.WriteString(fmt.Sprintf(" %s%d", n, h.census[n]))
	}
	buf.WriteString(fmt.Sprintf(" %d%% %d(%dK/s)",
		h.fullness, h.lastEvents, h.keps))
	return buf.String()
}

// Call outputs stat information to the terminal.
func (h *Stat) Call() {
	now := time.Now().UnixNano()

	h.m.Lock()
	defer h.m.Unlock()

	// reset for call
	h.names = make([]string, 0)
	h.census = make(map[string]int)
	pop := 0

	state := h.s.State()
	for _, a := range state.Sites {
		h.census[string(a.Type.Rune)]++
		pop++
	}

	for name := range h.census {
		h.names = append(h.names, name)
	}
	sort.Strings(h.names)

	// Write bounds dependent stats.
	width, height := state.Width, state.Height
	h.census[" "] = (width * height) - pop
	h.fullness = 100 * pop / (width * height)

	// Compute perf stats.
	events := 0
	if h.lastEvents <= state.Events {
		events = state.Events - h.lastEvents
		second := float64(now-h.lastClock) / float64(time.Second)
		h.keps = int(float64(events) / (1000 * second))
	}
	h.lastEvents = state.Events
	h.lastClock = now
}
