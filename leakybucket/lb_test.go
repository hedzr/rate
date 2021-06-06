package leakybucket_test

import (
	"fmt"
	"github.com/hedzr/cmdr/tool/randomizer"
	"github.com/hedzr/rate/leakybucket"
	"github.com/hedzr/rate/rateapi"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func BenchmarkRandInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rand.Int()
	}
}

func TestLeakyBucketLimiter(b *testing.T) {
	var counter int
	l := leakybucket.New(100, time.Second) // one req per 10ms
	defer l.Close()
	time.Sleep(300 * time.Millisecond)
	prev := time.Now()
	for i := 0; i < 20; i++ {
		now := l.TakeBlocked(1)
		counter++
		fmt.Println(i, now.Sub(prev))
		prev = now
		time.Sleep(1 * time.Millisecond)
	}
	b.Logf("%v requests allowed.", counter)
}

func TestLeakyBucketLimiterNonBlocked(b *testing.T) {
	var counter int
	l := leakybucket.New(100, time.Second) // one req per 10ms
	defer l.Close()
	time.Sleep(300 * time.Millisecond)
	prev := time.Now()
	for i := 0; i < 120; i++ {
		ok := l.Take(1)
		if !ok {
			b.Logf("#%d Take() returns not ok, counter: %v", i, l.Count())
			time.Sleep(50 * time.Millisecond)
		} else {
			//b.Logf("OK: #%d Take(), counter: %v", i, l.count)
			now := time.Now()
			counter++
			fmt.Println(i, now.Sub(prev))
			prev = now
			time.Sleep(1 * time.Millisecond)
		}
	}
	b.Logf("%v requests allowed.", counter)
}

func TestTB(b *testing.T) {
	var wg sync.WaitGroup
	var counter int64

	l := leakybucket.New(100, time.Second)
	defer l.Close()

	var r = randomizer.New()

	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			runner(b, l, &counter, r)
		}()
	}
	wg.Wait()
	b.Logf("%v requests allowed.", counter)
}

func runner(b *testing.T, l rateapi.Limiter, counter *int64, rand randomizer.Randomizer) {
	for i := 0; i < 100; i++ {
		ok := l.Take(1)
		if !ok {
			b.Logf("#%d Take() returns not ok, available: %v", i, l.Available())
			time.Sleep(100 * time.Millisecond)
		} else {
			//b.Logf("OK: #%d Take(), counter: %v", i, l.count)
			atomic.AddInt64(counter, 1)
			time.Sleep(time.Duration(rand.NextInRange(5, 15)) * time.Millisecond)
		}
	}
}
