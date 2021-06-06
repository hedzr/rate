package rate_test

import (
	"github.com/hedzr/rate"
	"github.com/hedzr/rate/rateapi"
	"testing"
	"time"
)

func TestRegister(t *testing.T) {
	var err error
	err = rate.Register("ab", func(maxCount int64, d time.Duration) rateapi.Limiter {
		return rate.New(rate.TokenBucket, maxCount, d)
	})
	if err != nil {
		t.Fatal(err)
	}
	err = rate.Register(string(rate.TokenBucket), func(maxCount int64, d time.Duration) rateapi.Limiter {
		return nil
	})
	if err == nil {
		t.Fatal("for reserved name the return should be an error")
	}

	if l := rate.New("ab", 100, time.Second); l != nil {
		defer l.Close()
		rate.Unregister("ab")
	} else {
		t.Fatal("expecting a limiter as the result returned")
	}
}
