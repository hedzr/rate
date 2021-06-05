package counter

import (
	"time"
)

func New(maxCount int64, d time.Duration) *counter {
	return &counter{
		int(maxCount),
		d,
		0,
		time.Now().UnixNano(),
	}
}

type counter struct {
	Maximal int
	Period  time.Duration
	count   int
	tick    int64 // in nanosecond
}

func (s *counter) take(count int) bool {
	if time.Now().UnixNano() > s.tick {
		// if timeout, reset counter regally at first
		s.count = 0
		s.tick = time.Now().Add(s.Period).UnixNano()
	}

	s.count += count            // it's acceptable in HPC scene
	return s.count <= s.Maximal // it's acceptable in HPC scene
}

func (s *counter) Take(count int) bool {
	ok := s.take(count)
	return ok
}

func (s *counter) TakeBlocked(count int) (requestAt time.Time) {
	requestAt = time.Now().UTC()
	ok := s.take(count)
	for !ok {
		time.Sleep(time.Duration((int64(s.Period) / int64(s.Maximal)) - (1000 - 1)))
		ok = s.take(count)
	}
	// time.Sleep(time.Duration(s.rate-int64(time.Now().Sub(requestAt))) - time.Millisecond)
	return
}

func (s *counter) Ticks() int64     { return s.tick }
func (s *counter) Count() int       { return s.count }
func (s *counter) Available() int64 { return int64(s.count) }
func (s *counter) Capacity() int64  { return int64(s.Maximal) }
func (s *counter) Close()           {}
