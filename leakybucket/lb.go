// Package leakybucket implements leaky-bucket algorithm
package leakybucket

import (
	"sync/atomic"
	"time"

	"github.com/hedzr/rate/pkg/logger"
	"github.com/hedzr/rate/rateapi"
)

// New make a new instance of limiter
func New(maxCount int64, d time.Duration) rateapi.Limiter {
	return (&leakyBucket{
		true,
		int64(maxCount),
		make(chan struct{}),
		int64(d) / int64(maxCount),
		time.Now().UnixNano(),
		0,
	}).start(d)
}

type leakyBucket struct {
	enabled     bool
	Maximal     int64
	exitCh      chan struct{}
	rate        int64
	refreshTime int64 // in nanoseconds
	count       int64
}

func (s *leakyBucket) Enabled() bool     { return s.enabled }
func (s *leakyBucket) SetEnabled(b bool) { s.enabled = b }
func (s *leakyBucket) Count() int64      { return atomic.LoadInt64(&s.count) }
func (s *leakyBucket) Available() int64  { return int64(atomic.LoadInt64(&s.count)) }
func (s *leakyBucket) Capacity() int64   { return int64(s.Maximal) }

func (s *leakyBucket) Close() {
	close(s.exitCh)
}

func (s *leakyBucket) start(d time.Duration) *leakyBucket {
	if s.rate < 1000 {
		logger.Errorf("the rate cannot be less than 1000us, it's %v", s.rate)
		return nil
	}

	// fmt.Printf("rate: %v\n", time.Duration(s.rate))

	// go s.looper(d)
	return s
}

//func (s *leakyBucket) looper(d time.Duration) {
//	// nothing to do
//}

func (s *leakyBucket) max(a, b int64) int64 {
	if a < b {
		return b
	}
	return a
}

func (s *leakyBucket) take(count int) (requestAt time.Time, ok bool) {
	requestAt = time.Now()

	rtm := atomic.LoadInt64(&s.refreshTime)
	cnt := atomic.LoadInt64(&s.count)
	atomic.StoreInt64(&s.count, s.max(0, cnt-(requestAt.UnixNano()-rtm)/s.rate*int64(count)))
	atomic.StoreInt64(&s.refreshTime, requestAt.UnixNano())

	cnt = atomic.LoadInt64(&s.count)
	if cnt < s.Maximal {
		atomic.AddInt64(&s.count, int64(count))
		ok = true
	}

	return
}

func (s *leakyBucket) Take(count int) (ok bool) {
	_, ok = s.take(count)
	return
}

func (s *leakyBucket) TakeBlocked(count int) (requestAt time.Time) {
	var ok bool
	requestAt, ok = s.take(count)
	for !ok {
		time.Sleep(time.Duration(s.rate - (1000 - 1)))
		_, ok = s.take(count)
	}
	time.Sleep(time.Duration(s.rate-int64(time.Now().Sub(requestAt))) - time.Millisecond)
	return
}
