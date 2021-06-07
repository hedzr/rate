package rateapi

import (
	"time"
)

// Limiter represents a rate-limiter generally
type Limiter interface {
	// Take assigns count of allows from a rate-limiter without blocking.
	Take(count int) bool
	// TakeBlocked assigns count of allows from a rate-limiter till requesting ok.
	TakeBlocked(count int) (requestAt time.Time)
	// Close is a Peripheral equivalent (see also hedzr/log and basics.Peripheral).
	// a rate limiter must be released safely at shutting down.
	Close()
	// Available returns a number for 'X-RateLimit-Remaining'
	Available() int64
	// Capacity returns a number for 'X-RateLimit-Limit'
	Capacity() int64

	// Enabled __
	Enabled() bool
	// SetEnabled __
	SetEnabled(b bool)
}

type endMeasurable interface {
	// Ticks represents the end-point in nanoseconds
	Ticks() int64
}

//type Countable interface {
//	Count() int
//}
//
//type Countable32 interface {
//	Count() int32
//}
//
//type Countable64 interface {
//	Count() int64
//}
