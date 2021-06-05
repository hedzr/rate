package tokenbucket_test

import (
	"fmt"
	"github.com/hedzr/rate/tokenbucket"
	"math/rand"
	"testing"
	"time"
)

func BenchmarkRandInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rand.Int()
	}
}

func TestTokenBucketLimiter(b *testing.T) {
	var counter int
	l := tokenbucket.New(100, time.Second) // one req per 10ms
	defer l.Close()
	time.Sleep(300 * time.Millisecond)
	prev := time.Now()
	for i := 0; i < 1200; i++ {
		now := l.TakeBlocked(1)
		fmt.Println(i, now.Sub(prev), l.Available())
		counter++
		prev = now
		//time.Sleep(1 * time.Millisecond)
	}
	b.Logf("%v requests allowed.", counter)
}

func TestTokenBucketLimiterNonBlocked(b *testing.T) {
	var counter int
	l := tokenbucket.New(100, time.Second) // one req per 10ms
	defer l.Close()
	time.Sleep(300 * time.Millisecond)
	for i := 0; i < 120; i++ {
		ok := l.Take(1)
		if !ok {
			b.Logf("#%d Take() returns not ok, counter: %v", i, l.Count())
			time.Sleep(100 * time.Millisecond)
		} else {
			//b.Logf("OK: #%d Take(), counter: %v", i, l.count)
			//time.Sleep(5 * time.Millisecond)
			counter++
		}
	}
	b.Logf("%v requests allowed.", counter)
}
