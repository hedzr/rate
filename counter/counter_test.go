package counter

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func BenchmarkRandInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rand.Int()
	}
}

func TestCounterLimiter(b *testing.T) {
	var counter int
	l := New(100, time.Second) // one req per 10ms
	defer l.Close()
	for i := 0; i < 120; i++ {
		ok := l.Take(1)
		if !ok {
			b.Logf("#%d Take() returns not ok, remained ticks: %vns, counter: %v", i, l.(interface{ Ticks() int64 }).Ticks()-time.Now().UnixNano(), l.(interface{ Count() int }).Count())
			time.Sleep(100 * time.Millisecond)
		} else {
			//time.Sleep(5 * time.Millisecond)
			counter++
		}
	}
	b.Logf("%v requests allowed.", counter)
}

func TestCounterLimiterBlocked(b *testing.T) {
	var counter int
	l := New(100, time.Second) // one req per 10ms
	defer l.Close()
	prev := time.Now()
	for i := 0; i < 120; i++ {
		now := l.TakeBlocked(1)
		fmt.Println(i, now.Sub(prev), l.Available())
		counter++
		prev = now
	}
	b.Logf("%v requests allowed.", counter)
	b.Log(l.Enabled(), l.Available(), l.Capacity())
	l.SetEnabled(false)
}
