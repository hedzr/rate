// Package rate provides a rate-limiter with a special algorithm.
package rate

import (
	"errors"
	"time"

	"github.com/hedzr/rate/counter"
	"github.com/hedzr/rate/leakybucket"
	"github.com/hedzr/rate/rateapi"
	"github.com/hedzr/rate/tokenbucket"
)

// Algorithm represents the rate limit algorithm specifically
type Algorithm string

const (
	// Counter algorithm
	Counter Algorithm = "counter"
	// LeakyBucket algorithm
	LeakyBucket Algorithm = "leaky-bucket"
	// TokenBucket algorithm
	TokenBucket Algorithm = "token-bucket"
)

// New returns a new instance of the rate limiter with certain a algorithm.
//
// a nil result means the algorithm of yours has not been registered, so you might:
//
// - use a right algorithm name such as rate.LeakyBucket, rate.TokenBucket
// - or register yours implement with rate.Register and assign it by algorithm name.
func New(algorithm Algorithm, maxCount int64, d time.Duration) rateapi.Limiter {
	if lfn, ok := knownLimiters[algorithm]; ok {
		return lfn(maxCount, d)
	}
	return nil
}

// CountOf extracts the Available tokens/rate-remains count from a rate-limiter
func CountOf(limiter rateapi.Limiter) int64 {
	if c, ok := limiter.(interface{ Count() int }); ok {
		return int64(c.Count())
	}
	if c, ok := limiter.(interface{ Count() int64 }); ok {
		return int64(c.Count())
	}
	if c, ok := limiter.(interface{ Count() int32 }); ok {
		return int64(c.Count())
	}
	return 0
}

// Register puts your generator into registry so it will be assign from New() in the future
func Register(algorithm string, generator func(maxCount int64, d time.Duration) rateapi.Limiter) error {
	switch Algorithm(algorithm) {
	case Counter, LeakyBucket, TokenBucket:
		return errors.New("reserved name found")
	}

	if _, ok := knownLimiters[Algorithm(algorithm)]; ok {
		return errors.New("name exists")
	}

	knownLimiters[Algorithm(algorithm)] = generator
	return nil
}

// Unregister deregister a limiter generator by algorithm name.
func Unregister(algorithm Algorithm) {
	if _, ok := knownLimiters[Algorithm(algorithm)]; ok {
		delete(knownLimiters, algorithm)
	}
}

func init() {
	knownLimiters = make(map[Algorithm]func(maxCount int64, d time.Duration) rateapi.Limiter)
	knownLimiters[Counter] = func(maxCount int64, d time.Duration) rateapi.Limiter {
		return counter.New(maxCount, d)
	}

	knownLimiters[LeakyBucket] = func(maxCount int64, d time.Duration) rateapi.Limiter {
		return leakybucket.New(maxCount, d)
	}

	knownLimiters[TokenBucket] = func(maxCount int64, d time.Duration) rateapi.Limiter {
		return tokenbucket.New(maxCount, d)
	}
}

// knownLimiters is a public registry to store the generators of a rate-limiter
var knownLimiters map[Algorithm]func(maxCount int64, d time.Duration) rateapi.Limiter
