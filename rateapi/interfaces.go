package rateapi

import (
	"time"
)

type Limiter interface {
	Take(count int) bool
	TakeBlocked(count int) (requestAt time.Time)
	Close()
	Available() int64
	Capacity() int64
}

type EndMeasurable interface {
	// Ticks represents the end-point in nanoseconds
	Ticks() int64
}

type Countable interface {
	Count() int
}

type Countable32 interface {
	Count() int32
}

type Countable64 interface {
	Count() int64
}

var KnownLimiters = make(map[string]func(maxCount int64, d time.Duration) Limiter)
