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

	err = rate.Register("ab", func(maxCount int64, d time.Duration) rateapi.Limiter {
		return rate.New(rate.TokenBucket, maxCount, d)
	})
	if err == nil {
		t.Fatal("cannot register again with a same name")
	}

	if l := rate.New("ab", 100, time.Second); l != nil {
		defer l.Close()
		rate.Unregister("ab")
	} else {
		t.Fatal("expecting a limiter as the result returned")
	}
}

func TestNilLimiter(t *testing.T) {
	l := rate.New("not-exists", 100, time.Second)
	if l != nil {
		t.Fatal("for reserved name the return should be an error")
	}
}

func TestCountOf(t *testing.T) {

	if l := rate.New(rate.TokenBucket, 100, time.Second); l != nil {
		defer l.Close()
		t.Logf("count: %v", rate.CountOf(l))
	} else {
		t.Fatal("expecting a limiter as the result returned")
	}

}

func TestNewLimiters(t *testing.T) {
	l := rate.New(rate.Counter, 100, time.Second)
	if l == nil {
		t.Fatal("New a limiter failed")
	}
	defer l.Close()

	l1 := rate.New(rate.LeakyBucket, 100, time.Second)
	if l1 == nil {
		t.Fatal("New a limiter failed")
	}
	defer l1.Close()

	l2 := rate.New(rate.TokenBucket, 100, time.Second)
	if l2 == nil {
		t.Fatal("New a limiter failed")
	}
	defer l2.Close()
}
