package counter

import (
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
			b.Logf("#%d Take() returns not ok, remained ticks: %vns, counter: %v", i, l.tick-time.Now().UnixNano(), l.count)
			time.Sleep(100 * time.Millisecond)
		} else {
			//time.Sleep(5 * time.Millisecond)
			counter++
		}
	}
	b.Logf("%v requests allowed.", counter)
}
