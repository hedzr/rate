package leakybucket_test

import (
	"fmt"
	"github.com/hedzr/rate/leakybucket"
	"math/rand"
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
	l := leakybucket.New(100, time.Second, false) // one req per 10ms
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
	l := leakybucket.New(100, time.Second, false) // one req per 10ms
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
