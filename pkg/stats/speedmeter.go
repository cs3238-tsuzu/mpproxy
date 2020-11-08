package stats

import (
	"sync/atomic"
	"time"
)

// Speedmeter is the interface to retrieve the current transfer speed
type Speedmeter interface {
	// Current returns the current transfer speed
	Current() int64
}

type speedmeter struct {
	tick  time.Duration
	bytes int64
	since time.Time
	bps   int64
	now   func() time.Time
}

var _ Speedmeter = &speedmeter{}

func newSpeedmeter(tick time.Duration) *speedmeter {
	return &speedmeter{
		tick: tick,
		now:  time.Now,
	}
}

func (s *speedmeter) update(bytes int64) {
	now := s.now()
	s.bytes += bytes

	if s.since == (time.Time{}) {
		s.since = now

		return
	}

	elapsed := now.Sub(s.since)
	if elapsed >= s.tick {
		atomic.StoreInt64(
			&s.bps,
			s.bytes*int64(time.Second)/int64(elapsed),
		)
		s.bytes = 0
		s.since = now
	}
}

// Current returns the current transfer speed
func (s *speedmeter) Current() int64 {
	return atomic.LoadInt64(&s.bps)
}
