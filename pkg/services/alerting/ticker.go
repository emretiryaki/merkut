package alerting


import (
	"time"
	"github.com/benbjohnson/clock"
)
type Ticker struct {
	C         chan time.Time
	clock     clock.Clock
	last      time.Time
	offset    time.Duration
	newOffset chan time.Duration
}

// NewTicker returns a ticker that ticks on second marks or very shortly after, and never drops ticks
func NewTicker(last time.Time, initialOffset time.Duration, c clock.Clock) *Ticker {
	t := &Ticker{
		C:         make(chan time.Time),
		clock:     c,
		last:      last,
		offset:    initialOffset,
		newOffset: make(chan time.Duration),
	}
	go t.run()
	return t
}

func (t *Ticker) updateOffset(offset time.Duration) {
	t.newOffset <- offset
}

func (t *Ticker) run() {
	for {
		next := t.last.Add(time.Duration(1) * time.Second)
		diff := t.clock.Now().Add(-t.offset).Sub(next)
		if diff >= 0 {
			t.C <- next
			t.last = next
			continue
		}
		// tick is too young. try again when ...
		select {
		case <-t.clock.After(-diff): // ...it'll definitely be old enough
		case offset := <-t.newOffset: // ...it might be old enough
			t.offset = offset
		}
	}
}
