package tokenbucket

import (
	"github.com/hedzr/log"
	"sync/atomic"
	"time"
)

func New(maxCount int64, d time.Duration) *tokenBucket {
	return (&tokenBucket{
		true,
		int32(maxCount),
		d,
		int64(d) / int64(maxCount),
		make(chan struct{}),
		int32(maxCount),
	}).start(d)
}

type tokenBucket struct {
	enabled bool
	Maximal int32
	period  time.Duration
	rate    int64
	exitCh  chan struct{}
	count   int32
}

func (s *tokenBucket) Enabled() bool     { return s.enabled }
func (s *tokenBucket) SetEnabled(b bool) { s.enabled = b }
func (s *tokenBucket) Count() int32      { return atomic.LoadInt32(&s.count) }
func (s *tokenBucket) Available() int64  { return int64(atomic.LoadInt32(&s.count)) }
func (s *tokenBucket) Capacity() int64   { return int64(s.Maximal) }

func (s *tokenBucket) Close() {
	close(s.exitCh)
}

func (s *tokenBucket) start(d time.Duration) *tokenBucket {
	if s.rate < 1000 {
		log.Errorf("the rate cannot be less than 1000us, it's %v", s.rate)
		return nil
	}

	go s.looper(d)
	return s
}

func (s *tokenBucket) looper(d time.Duration) {
	ticker := time.NewTicker(d / time.Duration(s.Maximal))
	// fmt.Printf("token building spped is: 1req/%v\n", d/time.Duration(s.Maximal))
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case <-s.exitCh:
			return
		case <-ticker.C:
			vn := atomic.AddInt32(&s.count, 1)
			if vn < s.Maximal {
				continue
			}

			vn %= s.Maximal
			if vn > 0 {
				atomic.StoreInt32(&s.count, s.Maximal)
			}
		}
	}
}

func (s *tokenBucket) take(count int) bool {
	if vn := atomic.AddInt32(&s.count, -1*int32(count)); vn >= 0 {
		return true
	}
	atomic.StoreInt32(&s.count, 0)
	return false
}

func (s *tokenBucket) Take(count int) bool {
	ok := s.take(count)
	return ok
}

func (s *tokenBucket) TakeBlocked(count int) (requestAt time.Time) {
	requestAt = time.Now().UTC()
	ok := s.take(count)
	for !ok {
		time.Sleep(time.Duration(s.rate - (1000 - 1)))
		ok = s.take(count)
	}
	// time.Sleep(time.Duration(s.rate-int64(time.Now().Sub(requestAt))) - time.Millisecond)
	return
}
