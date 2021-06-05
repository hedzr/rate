package rate

import (
	"github.com/hedzr/rate/counter"
	"github.com/hedzr/rate/leakybucket"
	"github.com/hedzr/rate/rateapi"
	"github.com/hedzr/rate/tokenbucket"
	"time"
)

func New(algorithm Algorithm, maxCount int64, d time.Duration) rateapi.Limiter {
	if lfn, ok := rateapi.KnownLimiters[string(algorithm)]; ok {
		return lfn(maxCount, d)
	}
	return nil
}

func CountOf(limiter rateapi.Limiter) int64 {
	if c, ok := limiter.(rateapi.Countable); ok {
		return int64(c.Count())
	} else if c, ok := limiter.(rateapi.Countable64); ok {
		return int64(c.Count())
	} else if c, ok := limiter.(rateapi.Countable32); ok {
		return int64(c.Count())
	}
	return 0
}

func init() {

	rateapi.KnownLimiters[string(Counter)] = func(maxCount int64, d time.Duration) rateapi.Limiter {
		return counter.New(maxCount, d)
	}

	rateapi.KnownLimiters[string(LeakyBucket)] = func(maxCount int64, d time.Duration) rateapi.Limiter {
		return leakybucket.New(maxCount, d)
	}
	//rateapi.KnownLimiters[string(LeakyBucketWithChannel)] = func(maxCount int64, d time.Duration) rateapi.Limiter {
	//	return leakybucket.New(maxCount, d, true)
	//}

	rateapi.KnownLimiters[string(TokenBucket)] = func(maxCount int64, d time.Duration) rateapi.Limiter {
		return tokenbucket.New(maxCount, d)
	}
}

type Algorithm string

const (
	Counter                Algorithm = "counter"
	LeakyBucket            Algorithm = "leaky-bucket"
	LeakyBucketWithChannel Algorithm = "leaky-bucket-with-out-channel"
	TokenBucket            Algorithm = "token-bucket"
)
