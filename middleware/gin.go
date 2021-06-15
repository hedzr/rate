package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hedzr/rate"
	"github.com/hedzr/rate/rateapi"
	"gopkg.in/hedzr/errors.v2"
	"time"
)

// ForGin builds a gin.HandlerFunc from a Config object
func ForGin(config *Config) gin.HandlerFunc {
	l := &exLimiter{
		enabled:   config.Enabled,
		algorithm: config.Algorithm,
		d:         config.Interval,
		capacity:  int64(config.MaxRequests),
		limiters:  make(map[string]rateapi.Limiter),
	}
	l.rateKeygen = l.buildKeyFunc(config)
	return l.Middleware()
}

// Middleware interface
type Middleware interface {
	Middleware() gin.HandlerFunc
}

// NewLimiterForGin builds a generic limiter wrapper which provides a Middleware() api for gin.Use.
func NewLimiterForGin(interval time.Duration, capacity int64, keyGen KeygenFunc) Middleware {
	limiters := make(map[string]rateapi.Limiter)
	return &exLimiter{
		true,
		string(rate.TokenBucket),
		interval,
		capacity,
		keyGen,
		limiters,
	}
}

func (r *exLimiter) buildKeyFunc(config *Config) KeygenFunc {
	return func(ctx *gin.Context) (string, error) {
		key := ctx.Request.Header.Get(config.HeaderKeyName)
		if key != "" {
			for _, k := range config.ExceptionKeys {
				if key == k {
					return passedBucketName, ErrRateLimitPassed
				}
			}
			return key, nil
		}
		msg := fmt.Sprintf("header key %q is missing", config.HeaderKeyName)
		ctx.JSON(403, gin.H{"code": 2901, "message": msg})
		return "", errors.New(msg)
	}
}

// KeygenFunc returns a functor to build a unique key name which can be used to associate with a limiter object
type KeygenFunc func(ctx *gin.Context) (uniqueKeyName string, err error)

type exLimiter struct {
	enabled    bool
	algorithm  string
	d          time.Duration
	capacity   int64
	rateKeygen KeygenFunc
	limiters   map[string]rateapi.Limiter
}

func (r *exLimiter) get(ctx *gin.Context) (rateapi.Limiter, error) {
	key, err := r.rateKeygen(ctx)

	if err != nil {
		if err != ErrRateLimitPassed {
			return nil, err
		}
		err, key = nil, passedBucketName
	}

	if limiter, existed := r.limiters[key]; existed {
		return limiter, nil
	}

	if key == passedBucketName {
		limiter := rate.New(rate.TokenBucket, int64(time.Minute/1000), time.Minute)
		r.limiters[key] = limiter
		return limiter, nil
	}

	limiter := rate.New(rate.TokenBucket, r.capacity, r.d)
	r.limiters[key] = limiter
	return limiter, nil
}

func (r *exLimiter) Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limiter, err := r.get(ctx)
		if err != nil {
			ctx.AbortWithError(429, err)
		} else if limiter != nil && limiter.Enabled() {
			if limiter.Take(1) {
				ctx.Writer.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", limiter.Available()))
				ctx.Writer.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", limiter.Capacity()))
				// ctx.Writer.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", limiter.Take(1)))
				ctx.Next()
			} else {
				err = errors.New("Too many requests")
				ctx.Writer.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", limiter.Available()))
				ctx.Writer.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", limiter.Capacity()))
				ctx.AbortWithError(429, err)
				// log.Warnf("[rate-limit] overflow: %v", e.Error())}
			}
		} else {
			ctx.Next()
		}
	}
}

//func (r *exLimiter) getRL(ctx gin.Context) rateapi.exLimiter {
//	if r.limiter == nil {
//		r.limiter = leakybucket.New(100, time.Second, false)
//	}
//	return r.limiter
//}
//
//func (r *exLimiter) SimpleMiddleware() gin.HandlerFunc {
//	return func(ctx *gin.Context) {
//		limiter := r.getRL(ctx)
//		if err != nil {
//			ctx.AbortWithError(429, err)
//		} else if limiter != nil {
//			if limiter.Take(1) {
//				ctx.Writer.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", limiter.Available()))
//				ctx.Writer.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", limiter.Capacity()))
//				// ctx.Writer.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", limiter.Take(1)))
//				ctx.Next()
//			} else {
//				err = errors.New("Too many requests")
//				ctx.Writer.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", limiter.Available()))
//				ctx.Writer.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", limiter.Capacity()))
//				ctx.AbortWithError(429, err)
//				// log.Warnf("[rate-limit] overflow: %v", e.Error())}
//			}
//		} else {
//			ctx.Next()
//		}
//	}
//}

// ErrRateLimitPassed identify a special state that an exception key was found.
// The limiter shouldn't be applied to the request which has been tagged with the exception key.
//
var ErrRateLimitPassed = errors.New("always passed up for exceptions")

const passedBucketName = "exceptions-met"
